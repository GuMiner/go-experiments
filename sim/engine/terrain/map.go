package terrain

import "go-experiments/sim/config"

type TerrainTexel struct {
	TerrainType TerrainType

	// Absolute height
	Height float32

	// Relative height for the given terrain type
	HeightPercent float32
}

type TerrainSubMap struct {
	Texels [][]TerrainTexel
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
		Texels: make([][]TerrainTexel, regionSize*regionSize)}

	for i := 0; i < regionSize; i++ {
		terrainSubMap.Texels[i] = make([]TerrainTexel, regionSize)
	}

	heights := Generate(regionSize, regionSize, x*regionSize, y*regionSize)
	for i := 0; i < regionSize; i++ {
		for j := 0; j < regionSize; j++ {
			height := heights[i+j*regionSize]
			terrainType, percent := GetTerrainType(height)
			terrainSubMap.Texels[i][j] = TerrainTexel{
				TerrainType:   terrainType,
				Height:        height,
				HeightPercent: percent}
		}
	}

	return &terrainSubMap
}

func (t *TerrainMap) GetOrAddRegion(x, y int) *TerrainSubMap {
	if _, ok := t.SubMaps[x]; !ok {
		t.SubMaps[x] = make(map[int]*TerrainSubMap)
	}

	if _, ok := t.SubMaps[x][y]; !ok {
		t.SubMaps[x][y] = NewTerrainSubMap(x, y)
	}

	return t.SubMaps[x][y]
}
