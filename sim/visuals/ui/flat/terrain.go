package flat

import (
	"go-experiments/sim/config"
	"go-experiments/sim/engine/terrain"

	"github.com/go-gl/mathgl/mgl32"
)

// Given a height, returns the terrain color and percentage within that level
func GetTerrainColor(height float32) (mgl32.Vec3, float32) {
	terrainType, percent := terrain.GetTerrainType(height)

	switch terrainType {
	case terrain.Water:
		return config.Config.Ui.TerrainUi.WaterColor.ToVec3(), percent
	case terrain.Sand:
		return config.Config.Ui.TerrainUi.SandColor.ToVec3(), percent
	case terrain.Grass:
		return config.Config.Ui.TerrainUi.GrassColor.ToVec3(), percent
	case terrain.Hills:
		return config.Config.Ui.TerrainUi.HillColor.ToVec3(), percent
	case terrain.Rocks:
		return config.Config.Ui.TerrainUi.RockColor.ToVec3(), percent
	default:
		return config.Config.Ui.TerrainUi.SnowColor.ToVec3(), percent
	}
}
