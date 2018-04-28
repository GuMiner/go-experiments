package terrain

import "go-experiments/sim/config"

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

// Given a height, returns the terrain type and percentage within that level
func GetTerrainType(height float32) (TerrainType, float32) {
	if height < config.Config.Terrain.WaterLevel {
		percent := getLevelPercent(height, 0, config.Config.Terrain.WaterLevel)
		return Water, percent
	} else if height < config.Config.Terrain.SandLevel {
		percent := getLevelPercent(height, config.Config.Terrain.WaterLevel, config.Config.Terrain.SandLevel)
		return Sand, percent
	} else if height < config.Config.Terrain.GrassLevel {
		percent := getLevelPercent(height, config.Config.Terrain.SandLevel, config.Config.Terrain.GrassLevel)
		return Grass, percent
	} else if height < config.Config.Terrain.HillLevel {
		percent := getLevelPercent(height, config.Config.Terrain.GrassLevel, config.Config.Terrain.HillLevel)
		return Hills, percent
	} else if height < config.Config.Terrain.RockLevel {
		percent := getLevelPercent(height, config.Config.Terrain.HillLevel, config.Config.Terrain.RockLevel)
		return Rocks, percent
	}

	percent := getLevelPercent(height, config.Config.Terrain.RockLevel, config.Config.Terrain.SnowLevel)
	return Snow, percent
}

func getLevelPercent(height float32, min float32, max float32) float32 {
	return (height - min) / (max - min)
}
