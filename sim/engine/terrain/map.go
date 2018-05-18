package terrain

import (
	"fmt"
	"go-experiments/common/commonmath"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/subtile"

	"github.com/go-gl/mathgl/mgl32"
)

type TerrainTexel struct {
	TerrainType TerrainType

	// Absolute height
	Height float32

	// Relative height for the given terrain type
	HeightPercent float32
}

func (t *TerrainTexel) Normalize() {
	t.Height = commonMath.MinFloat32(1, commonMath.MaxFloat32(0, t.Height))
	t.TerrainType, t.HeightPercent = GetTerrainType(t.Height)
}

type TerrainSubMap struct {
	Texels [][]TerrainTexel

	// Flags for knowing when asynchronous generation is complete
	generated                bool
	generationCompleteSignal chan bool

	// Flag to indicate that the terrain is dirty and must be re-rendered
	Dirty bool
}

type TerrainMap struct {
	SubMaps map[int]map[int]*TerrainSubMap
}

func NewTerrainMap() *TerrainMap {
	terrainMap := TerrainMap{
		SubMaps: make(map[int]map[int]*TerrainSubMap)}

	return &terrainMap
}

func NewTerrainSubMap(x, y int) *TerrainSubMap {
	regionSize := config.Config.Terrain.RegionSize

	terrainSubMap := TerrainSubMap{
		Texels:                   make([][]TerrainTexel, regionSize*regionSize),
		generated:                false,
		generationCompleteSignal: make(chan bool)}

	for i := 0; i < regionSize; i++ {
		terrainSubMap.Texels[i] = make([]TerrainTexel, regionSize)
	}

	return &terrainSubMap
}

func (t *TerrainSubMap) GenerateSubMap(x, y int) {
	regionSize := config.Config.Terrain.RegionSize
	heights := Generate(regionSize, regionSize, x*regionSize, y*regionSize)
	for i := 0; i < regionSize; i++ {
		for j := 0; j < regionSize; j++ {
			height := heights[i+j*regionSize]
			t.Texels[i][j] = TerrainTexel{Height: height}
			t.Texels[i][j].Normalize()
		}
	}

	t.generated = true

	fmt.Printf("Generated sub map terrain for [%v, %v]\n", x, y)
	t.generationCompleteSignal <- true
	close(t.generationCompleteSignal)

	fmt.Printf("  Terrain sub map [%v, %v] consumed.\n", x, y)
}

// Adds a region to the map, without waiting for its generation to complete.
func (t *TerrainMap) AddRegionIfMissing(x, y int) {
	if _, ok := t.SubMaps[x]; !ok {
		t.SubMaps[x] = make(map[int]*TerrainSubMap)
	}

	if _, ok := t.SubMaps[x][y]; !ok {
		t.SubMaps[x][y] = NewTerrainSubMap(x, y)
		go t.SubMaps[x][y].GenerateSubMap(x, y)
	}
}

func (t *TerrainMap) GetOrAddRegion(x, y int) *TerrainSubMap {
	t.AddRegionIfMissing(x, y)

	// If added but not generated, a generation thread is running
	// Wait for it before continuing.
	if !t.SubMaps[x][y].generated {
		_ = <-t.SubMaps[x][y].generationCompleteSignal
	}

	return t.SubMaps[x][y]
}

func (t *TerrainMap) ValidateGroundLocation(reg commonMath.Region) bool {

	iterate := func(x, y int) bool {
		pos := mgl32.Vec2{float32(x), float32(y)}
		texel, _ := t.getTexel(pos)

		return texel.TerrainType == Water
	}

	return !reg.IterateIntWithEarlyExit(iterate)
}

func (t *TerrainMap) Flatten(region commonMath.Region, amount float32) {
	t.performRegionBasedUpdate(region, amount, flatten)
}

func (t *TerrainMap) Sharpen(region commonMath.Region, amount float32) {
	t.performRegionBasedUpdate(region, amount, sharpen)
}

func (t *TerrainMap) Hills(region commonMath.Region, amount float32) {
	t.performRegionBasedUpdate(region, amount, hills)
}

func (t *TerrainMap) Valleys(region commonMath.Region, amount float32) {
	t.performRegionBasedUpdate(region, amount, valleys)
}

func (t *TerrainMap) performRegionBasedUpdate(region commonMath.Region, amount float32, update func(mgl32.Vec2, mgl32.Vec2, *TerrainTexel, float32, float32, float32)) {
	centerTexel, _ := t.getTexel(region.Position)
	centralHeight := centerTexel.Height

	region.IterateIntWithEarlyExit(func(x, y int) bool {
		modifiedPos := mgl32.Vec2{float32(x) + 0.5, float32(y) + 0.5}
		texel, texelRegion := t.getTexel(modifiedPos)

		update(region.Position, modifiedPos, texel, centralHeight, amount, region.Scale)
		texelRegion.Dirty = true

		// Never early exit
		return false
	})
}

// Average, moving parts that are farther away closer in faster.
func flatten(centerPosition, texelPosition mgl32.Vec2, texel *TerrainTexel, centerHeight, amount, regionSize float32) {
	heightDifference := texel.Height - centerHeight
	texel.Height = texel.Height - heightDifference*amount
	texel.Normalize()
}

// Reverse average, moving parts that are farther away further faster.
func sharpen(centerPosition, texelPosition mgl32.Vec2, texel *TerrainTexel, centerHeight, amount, regionSize float32) {
	heightDifference := texel.Height - centerHeight
	texel.Height = texel.Height + heightDifference*amount
	texel.Normalize()
}

// Makes hills, pushing pixels near the center position upwards,
func hills(centerPosition, texelPosition mgl32.Vec2, texel *TerrainTexel, centerHeight, amount, regionSize float32) {
	distanceFactor := 1.0 - centerPosition.Sub(texelPosition).Len()/regionSize

	texel.Height = texel.Height + amount*distanceFactor
	texel.Normalize()
}

func valleys(centerPosition, texelPosition mgl32.Vec2, texel *TerrainTexel, centerHeight, amount, regionSize float32) {
	distanceFactor := 1.0 - centerPosition.Sub(texelPosition).Len()/regionSize

	texel.Height = texel.Height - amount*distanceFactor
	texel.Normalize()
}

func (t *TerrainMap) getTexel(pos mgl32.Vec2) (*TerrainTexel, *TerrainSubMap) {
	regionX, regionY := subtile.GetRegionIndices(pos, config.Config.Terrain.RegionSize)
	region := t.GetOrAddRegion(regionX, regionY)

	localX, localY := subtile.GetLocalIndices(pos, regionX, regionY, config.Config.Terrain.RegionSize)
	return &region.Texels[localX][localY], region
}
