package engine

import (
	"go-experiments/common/math"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/element"
	"go-experiments/sim/engine/power"
	"go-experiments/sim/engine/terrain"
	"go-experiments/sim/input/editorEngine"
	"go-experiments/voxelli/utils"

	"github.com/go-gl/mathgl/mgl32"
)

type Engine struct {
	terrainMap    *terrain.TerrainMap
	elementFinder *element.ElementFinder
	powerPlants   *power.PowerPlants

	isMousePressed bool
	lastBoardPos   mgl32.Vec2
}

func NewEngine() *Engine {
	terrain.Init(config.Config.Terrain.Generation.Seed)

	engine := Engine{
		terrainMap:     terrain.NewTerrainMap(),
		elementFinder:  element.NewElementFinder(),
		powerPlants:    power.NewPowerPlants(),
		isMousePressed: false}

	return &engine
}

func (e *Engine) MousePress(pos mgl32.Vec2, engineState editorEngine.State) {
	e.isMousePressed = true
	e.lastBoardPos = pos
	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
		plantSize := power.Small // TODO: Configurable

		// Ensure we only put power plants on valid ground.
		_, size := power.GetPowerOutputAndSize(plantType, plantSize)

		// TODO: Get region from the item, not by hardcoding it here and hypothetically here.
		region := commonMath.Region{
			RegionType:  commonMath.SquareRegion,
			Scale:       float32(size),
			Orientation: 0,
			Position:    pos}

		// TODO: We really want 'any intersecting objects'. This is good for starters and infrastructure, but not much else.
		anyNearbyObjects := e.elementFinder.AnyInRange(pos, float32(size))

		if !anyNearbyObjects && e.terrainMap.ValidateGroundLocation(region) {
			element := e.powerPlants.Add(pos, plantType, plantSize)
			e.elementFinder.Add(element)
		}
	}
}

func (e *Engine) MouseMoved(pos mgl32.Vec2) {
	e.lastBoardPos = pos
}

func (e *Engine) MouseRelease(pos mgl32.Vec2, engineState editorEngine.State) {
	e.isMousePressed = false
	e.lastBoardPos = pos
}

func (e *Engine) Step(stepAmount float32, engineState editorEngine.State) {
	if engineState.Mode == editorEngine.Draw {
		if e.isMousePressed {
			switch engineState.InDrawMode {
			case editorEngine.TerrainFlatten:
				// TODO: We'll need to modify this to pass in a range...
				e.terrainMap.Flatten(e.lastBoardPos, stepAmount)
			default:
				break
			}
		}
	}
}

// Returns true if there is a hypothetical region that should currently be displayed, false otherwise.
func (e *Engine) HasHypotheticalRegion(pos mgl32.Vec2, engineState editorEngine.State) bool {
	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		return true
	} else if engineState.Mode == editorEngine.Draw {
		return true
	}

	return false
}

// Gets the hypothetical region that an action will happen to when the mouse is released.
func (e *Engine) GetHypotheticalRegion(pos mgl32.Vec2, engineState editorEngine.State) (isValid bool, region commonMath.Region) {
	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		// TODO this really shouldn't be duplicated with the above.
		plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
		plantSize := power.Small // TODO: Configurable

		// Ensure we only put power plants on valid ground.
		_, size := power.GetPowerOutputAndSize(plantType, plantSize)
		region := commonMath.Region{
			RegionType:  commonMath.SquareRegion,
			Scale:       float32(size),
			Orientation: 0,
			Position:    pos}

		anyNearbyObjects := e.elementFinder.AnyInRange(pos, float32(size))
		return !anyNearbyObjects && e.terrainMap.ValidateGroundLocation(region), region
	} else if engineState.Mode == editorEngine.Draw {
		return true, commonMath.Region{
			RegionType:  commonMath.CircleRegion,
			Scale:       30.0, // TODO Make this configurable by reading the editor engine state.
			Orientation: 0,    // It's a circle, so this doesn't really matter.
			Position:    pos}
	}

	return false, commonMath.Region{}
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
