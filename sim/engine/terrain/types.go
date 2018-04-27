package terrain

type TerrainType int

// Ordering is important! Reordering these changes their numerical values
const (
	Water TerrainType = iota
	Sand
	Grass
	Hills
	Rocks
	Snow
)
