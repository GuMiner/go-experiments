package engine

import (
	"go-experiments/sim/config"
	"go-experiments/sim/engine/power"
	"go-experiments/sim/engine/terrain"
	"go-experiments/voxelli/utils"
)

type Engine struct {
	terrainMap  *terrain.TerrainMap
	powerPlants *power.PowerPlants
}

func NewEngine() *Engine {
	terrain.Init(config.Config.Terrain.Generation.Seed)

	engine := Engine{
		terrainMap:  terrain.NewTerrainMap(),
		powerPlants: power.NewPowerPlants()}

	return &engine
}

func (e *Engine) ABC() {

}

// Update methods based on UI
func (e *Engine) PrecacheRegions(regions []utils.IntVec2) {
	for _, region := range regions {
		e.terrainMap.AddRegionIfMissing(region.X(), region.Y())
	}
}

// Data retrieval for drawing
func (e *Engine) GetRegionMap(region utils.IntVec2) *terrain.TerrainSubMap {
	return e.terrainMap.GetOrAddRegion(region.X(), region.Y())
}
