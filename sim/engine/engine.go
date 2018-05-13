package engine

import (
	"go-experiments/common/math"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/power"
	"go-experiments/sim/engine/terrain"
	"go-experiments/sim/input/editorEngine"
	"go-experiments/voxelli/utils"

	"github.com/go-gl/mathgl/mgl32"
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

func (e *Engine) MousePress(pos mgl32.Vec2, engineState editorEngine.State) {
	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
		plantSize := power.Small // TODO: Configurable

		// Ensure we only put power plants on valid ground.
		_, size := power.GetPowerOutputAndSize(plantType, plantSize)
		if e.terrainMap.ValidateGroundLocation(pos, size) {
			e.powerPlants.Add(pos, plantType, plantSize)
		}
	}
}

// Returns true if there is a hypothetical region that should currently be displayed, false otherwise.
func (e *Engine) HasHypotheticalRegion(pos mgl32.Vec2, engineState editorEngine.State) bool {
	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		return true
	}

	return false
}

// Gets the hypothetical region that an action will happen to when the mouse is released.
func (e *Engine) GetHypotheticalRegion(pos mgl32.Vec2, engineState editorEngine.State) (isValid bool, region commonMath.Region) {

	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
		plantSize := power.Small // TODO: Configurable

		// Ensure we only put power plants on valid ground.
		_, size := power.GetPowerOutputAndSize(plantType, plantSize)
		region := commonMath.Region{
			RegionType: commonMath.SquareRegion,
			Size:       float32(size)}

		// TODO, we also need to validate that there is not another plant or other structure in the way.

		return e.terrainMap.ValidateGroundLocation(pos, size), region
	}

	return false, commonMath.Region{}
}

func (e *Engine) MouseRelease(pos mgl32.Vec2, engineState editorEngine.State) {

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
