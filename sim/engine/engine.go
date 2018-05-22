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
	powerLineState  *PowerLineEditState
	snapElements    SnapElements

	Hypotheticals HypotheticalActions
}

func NewEngine() *Engine {
	terrain.Init(config.Config.Terrain.Generation.Seed)

	engine := Engine{
		terrainMap:      terrain.NewTerrainMap(),
		elementFinder:   element.NewElementFinder(),
		powerGrid:       power.NewPowerGrid(),
		isMousePressed:  false,
		actionPerformed: false,
		powerLineState:  NewPowerLineEditState(),
		snapElements:    NewSnapElements(),

		Hypotheticals: NewHypotheticalActions()}
	return &engine
}

func (e *Engine) addPowerPlantIfValid() {
	intesectsWithElement := e.elementFinder.IntersectsWithElement(e.getEffectivePosition(), e.Hypotheticals.Regions[0].Region.Scale)

	if !intesectsWithElement {
		isGroundValid := e.terrainMap.ValidateGroundLocation(e.Hypotheticals.Regions[0].Region)
		if isGroundValid {
			plantType := power.GetPlantType(editorEngine.EngineState.InPowerPlantAddMode)
			plantSize := power.Small // TODO: Configurable

			element := e.powerGrid.Add(e.getEffectivePosition(), plantType, plantSize)
			e.elementFinder.Add(element)
		}
	}
}

func (e *Engine) updatePowerLineState() {
	// If this is the first press, we associate it with the first location of the powerline.
	if !e.powerLineState.hasFirstNode {
		e.powerLineState.firstNode = e.getEffectivePosition()
		e.powerLineState.hasFirstNode = true
		e.powerLineState.firstNodeElement = e.getEffectivePowerGridElement()
	} else {
		// TODO: Configurable capacity
		line := e.powerGrid.AddLine(e.powerLineState.firstNode,
			e.getEffectivePosition(), 1000,
			e.powerLineState.firstNodeElement, e.getEffectivePowerGridElement())
		if line != nil {
			e.elementFinder.Add(line)

			e.powerLineState.firstNode = e.getEffectivePosition()
			e.powerLineState.firstNodeElement = line.GetSnapNodeElement(1)
		}
	}
}

func (e *Engine) getEffectivePosition() mgl32.Vec2 {
	if e.snapElements.snappedNode != nil {
		return e.snapElements.snappedNode.Element.GetSnapNodes()[e.snapElements.snappedNode.SnapNodeIdx]
	}

	if e.snapElements.snappedGridPos != nil {
		return *e.snapElements.snappedGridPos
	}

	return e.lastBoardPos
}

// TODO: Rename, element is too generic...
func (e *Engine) getEffectivePowerGridElement() int {
	node := e.snapElements.snappedNode
	if node != nil {
		// TODO: New interface for power elements?
		if line, ok := node.Element.(*power.PowerLine); ok {
			return line.GetSnapNodeElement(node.SnapNodeIdx)
		}

		if powerPlant, ok := node.Element.(*power.PowerPlant); ok {
			return powerPlant.GetSnapElement()
		}
	}

	// No grid element association.
	return -1
}

func (e *Engine) MousePress(pos mgl32.Vec2, engineState editorEngine.State) {
	e.isMousePressed = true
	e.lastBoardPos = pos
	if !e.actionPerformed {
		if engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerPlant {
			e.addPowerPlantIfValid()
		}
	}
}

func (e *Engine) MouseMoved(pos mgl32.Vec2, engineState editorEngine.State) {
	e.lastBoardPos = pos

	e.snapElements.ComputeSnappedSnapElements(e.lastBoardPos, e.elementFinder)
	e.powerLineState.EnterOrExitEditMode(&engineState)
}

func (e *Engine) MouseRelease(pos mgl32.Vec2, engineState editorEngine.State) {
	e.isMousePressed = false
	e.actionPerformed = false
	e.lastBoardPos = pos

	if e.powerLineState.InPowerLineState(&engineState) {
		e.updatePowerLineState()
	}
}

// Cancels the state of any multi-step operation, resetting it back to the start.
func (e *Engine) CancelState(engineState editorEngine.State) {
	if e.powerLineState.InPowerLineState(&engineState) {
		e.powerLineState.Reset()
	}
}

func (e *Engine) applyStepDraw(stepAmount float32, engineState *editorEngine.State) {
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

// Performs operations that are performed as steps with time.
func (e *Engine) Step(stepAmount float32, engineState editorEngine.State) {
	if engineState.Mode == editorEngine.Draw && e.isMousePressed {
		e.applyStepDraw(stepAmount, &engineState)
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

func (e *Engine) GetSnapElements() *SnapElements {
	return &e.snapElements
}
