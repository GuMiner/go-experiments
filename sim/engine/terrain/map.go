package terrain

import (
	"fmt"
	"go-experiments/common/math"
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

	generated                bool
	generationCompleteSignal chan bool
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
		regionX, regionY := subtile.GetRegionIndices(pos, config.Config.Terrain.RegionSize)
		region := t.GetOrAddRegion(regionX, regionY)

		localX, localY := subtile.GetLocalIndices(pos, regionX, regionY, config.Config.Terrain.RegionSize)
		centralTexel := region.Texels[localX][localY]

		return centralTexel.TerrainType == Water
	}

	return !reg.IterateIntWithEarlyExit(iterate)
}
