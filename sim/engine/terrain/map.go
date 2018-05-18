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
			terrainType, percent := GetTerrainType(height)
			t.Texels[i][j] = TerrainTexel{
				TerrainType:   terrainType,
				Height:        height,
				HeightPercent: percent}
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

func (t *TerrainMap) Flatten(pos mgl32.Vec2, amount float32) {
	// TODO: This should be a circle, not a square, and we should not be hardcoding the size!
	halfSize := 15

	// TODO: Should also be configurable
	amount = amount * 0.1

	centerTexel, _ := t.getTexel(pos)
	centralHeight := centerTexel.Height

	// TODO: This works, but there are edge effects. There should be a getTexel overload that works on the integer level.
	for i := int(pos.X()) - halfSize; i <= int(pos.X())+halfSize; i++ {
		for j := int(pos.Y()) - halfSize; j <= int(pos.Y())+halfSize; j++ {

			modifiedPos := mgl32.Vec2{float32(i), float32(j)}
			texel, region := t.getTexel(modifiedPos)

			// Average, moving parts that are farther away closer in faster. Effectively, flattening.
			heightDifference := texel.Height - centralHeight
			texel.Height = texel.Height - heightDifference*amount
			region.Dirty = true
		}
	}
}

func (t *TerrainMap) getTexel(pos mgl32.Vec2) (*TerrainTexel, *TerrainSubMap) {
	regionX, regionY := subtile.GetRegionIndices(pos, config.Config.Terrain.RegionSize)
	region := t.GetOrAddRegion(regionX, regionY)

	localX, localY := subtile.GetLocalIndices(pos, regionX, regionY, config.Config.Terrain.RegionSize)
	return &region.Texels[localX][localY], region
}
