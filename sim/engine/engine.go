package engine

import (
	"go-experiments/common/commonmath"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/element"
	"go-experiments/sim/engine/power"
	"go-experiments/sim/engine/terrain"
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/mathgl/mgl32"
)

type Engine struct {
	terrainMap    *terrain.TerrainMap
	elementFinder *element.ElementFinder
	powerGrid     *power.PowerGrid

	isMousePressed  bool
	actionPerformed bool
	lastBoardPos    mgl32.Vec2

	Hypotheticals HypotheticalActions
}

func NewEngine() *Engine {
	terrain.Init(config.Config.Terrain.Generation.Seed)

	engine := Engine{
		terrainMap:      terrain.NewTerrainMap(),
		elementFinder:   element.NewElementFinder(),
		powerGrid:       power.NewPowerGrid(),
		isMousePressed:  false,
		actionPerformed: false}

	engine.Hypotheticals.Reset()
	return &engine
}

func (e *Engine) MousePress(pos mgl32.Vec2, engineState editorEngine.State) {
	e.isMousePressed = true
	e.lastBoardPos = pos
	if !e.actionPerformed {
		if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
			// TODO: We really want 'any intersecting objects'. This is good for starters and infrastructure, but not much else.
			anyNearbyObjects := e.elementFinder.AnyInRange(pos, e.Hypotheticals.Regions[0].Region.Scale)

			if !anyNearbyObjects {
				isGroundValid := e.terrainMap.ValidateGroundLocation(e.Hypotheticals.Regions[0].Region)
				if isGroundValid {
					plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
					plantSize := power.Small // TODO: Configurable

					element := e.powerGrid.Add(pos, plantType, plantSize)
					e.elementFinder.Add(element)
				}
			}
		}
	}
}

func (e *Engine) MouseMoved(pos mgl32.Vec2) {
	e.lastBoardPos = pos
}

func (e *Engine) MouseRelease(pos mgl32.Vec2, engineState editorEngine.State) {
	e.isMousePressed = false
	e.actionPerformed = false
	e.lastBoardPos = pos
}

func (e *Engine) Step(stepAmount float32, engineState editorEngine.State) {
	if engineState.Mode == editorEngine.Draw {
		if e.isMousePressed {
			region := e.Hypotheticals.Regions[0].Region
			stepFactor := 0.1 * stepAmount

			switch engineState.InDrawMode {
			case editorEngine.TerrainFlatten:
				e.terrainMap.Flatten(region, stepFactor)
			case editorEngine.TerrainSharpen:
				e.terrainMap.Sharpen(region, stepFactor)
			case editorEngine.TerrainHills:
				e.terrainMap.Hills(region, stepFactor)
			case editorEngine.TerrainValleys:
				e.terrainMap.Valleys(region, stepFactor)
			default:
				break
			}
		}
	}
}

// Returns true if there is a hypothetical region that should currently be displayed, false otherwise.
func (e *Engine) ComputeHypotheticalRegion(engineState editorEngine.State) {
	if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
		e.Hypotheticals.Regions = make([]HypotheticalRegion, 1)

		plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
		plantSize := power.Small // TODO: Configurable

		// Ensure we only put power plants on valid ground.
		_, size := power.GetPowerOutputAndSize(plantType, plantSize)
		region := commonMath.Region{
			RegionType:  commonMath.SquareRegion,
			Scale:       float32(size),
			Orientation: 0,
			Position:    e.lastBoardPos}

		anyNearbyObjects := e.elementFinder.AnyInRange(e.lastBoardPos, region.Scale)
		var color mgl32.Vec3
		if !anyNearbyObjects && e.terrainMap.ValidateGroundLocation(region) {
			color = mgl32.Vec3{0, 1, 0}
		} else {
			color = mgl32.Vec3{1, 0, 0}
		}

		e.Hypotheticals.Regions[0] = HypotheticalRegion{
			Region: region,
			Color:  color}
	} else if engineState.Mode == editorEngine.Draw {
		e.Hypotheticals.Regions = []HypotheticalRegion{
			HypotheticalRegion{
				Color: mgl32.Vec3{0.0, 1.0, 1.0},
				Region: commonMath.Region{
					RegionType:  commonMath.CircleRegion,
					Scale:       30.0, // TODO Make this configurable by reading the editor engine state.
					Orientation: 0,
					Position:    e.lastBoardPos}}}
	} else {
		e.Hypotheticals.Reset()
	}
}

// Update methods based on UI
func (e *Engine) PrecacheRegions(regions []commonMath.IntVec2) {
	for _, region := range regions {
		e.terrainMap.AddRegionIfMissing(region.X(), region.Y())
	}
}

// Data retrieval for drawing
func (e *Engine) GetRegionMap(region commonMath.IntVec2) *terrain.TerrainSubMap {
	return e.terrainMap.GetOrAddRegion(region.X(), region.Y())
}

func (e *Engine) GetPowerGrid() *power.PowerGrid {
	return e.powerGrid
}

func (e *Engine) GetElementFinder() *element.ElementFinder {
	return e.elementFinder
}
