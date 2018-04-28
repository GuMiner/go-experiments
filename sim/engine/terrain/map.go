package terrain

type TerrainTexel struct {
	TerrainType TerrainType
	Height      float32
}

type TerrainSubMap struct {
	Texels [][]TerrainTexel
}

type TerrainMap struct {
	SubMaps map[int]map[int]TerrainSubMap
}
