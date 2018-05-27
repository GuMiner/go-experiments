package engine

import (
	"go-experiments/common/commonmath"
	"go-experiments/sim/engine/power"
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/mathgl/mgl32"
)

type HypotheticalRegion struct {
	Color  mgl32.Vec3
	Region commonMath.Region
}

type HypotheticalLine struct {
	Color mgl32.Vec3
	Line  [2]mgl32.Vec2
}

// Defines hypothetical regions for drawing and actions
type HypotheticalActions struct {
	Regions []HypotheticalRegion
	Lines   []HypotheticalLine
}

func NewHypotheticalActions() HypotheticalActions {
	h := HypotheticalActions{}
	h.Reset()
	return h
}

func (h *HypotheticalActions) Reset() {
	h.Regions = []HypotheticalRegion{}
	h.Lines = []HypotheticalLine{}
}

func (h *HypotheticalActions) setSingleRegion(region HypotheticalRegion) {
	h.Reset()
	h.Regions = []HypotheticalRegion{region}
}

func (e *HypotheticalActions) computePowerPlantHypotheticalRegion(n *Engine) {
	plantType := power.GetPlantType(editorEngine.EngineState.ItemSubSelection)
	plantSize := power.Small // TODO: Configurable

	// Ensure we only put power plants on valid ground.
	_, size := power.GetPowerOutputAndSize(plantType, plantSize)
	region := commonMath.Region{
		RegionType:  commonMath.SquareRegion,
		Scale:       float32(size),
		Orientation: 0,
		Position:    n.getEffectivePosition()}

	anyNearbyObjects := n.elementFinder.IntersectsWithElement(n.lastBoardPos, region.Scale)
	var color mgl32.Vec3
	if !anyNearbyObjects && n.terrainMap.ValidateGroundLocation(region) {
		color = mgl32.Vec3{0, 1, 0}
	} else {
		color = mgl32.Vec3{1, 0, 0}
	}

	e.setSingleRegion(HypotheticalRegion{Region: region, Color: color})
}

func (e *HypotheticalActions) computePowerLineHypotheticalRegion(n *Engine) {
	if !n.powerLineState.hasFirstNode {
		// Draw a generic powerline icon.
		e.setSingleRegion(
			HypotheticalRegion{
				Color: mgl32.Vec3{1.0, 0.0, 1.0},
				Region: commonMath.Region{
					RegionType:  commonMath.CircleRegion,
					Scale:       40.0, // TODO Make this configurable by reading the editor engine state.
					Orientation: 0,
					Position:    n.getEffectivePosition()}})
	} else {
		e.Reset()
		e.Lines = []HypotheticalLine{
			HypotheticalLine{
				Color: mgl32.Vec3{1.0, 0.0, 1.0},
				Line: [2]mgl32.Vec2{
					n.getEffectivePosition(),
					n.powerLineState.firstNode}}}
	}
}

func (e *HypotheticalActions) computeRoadLineHypotheticalRegion(n *Engine) {
	if !n.roadLineState.hasFirstNode {
		// Draw a generic road line icon.
		e.setSingleRegion(
			HypotheticalRegion{
				Color: mgl32.Vec3{1.0, 1.0, 0.0},
				Region: commonMath.Region{
					RegionType:  commonMath.CircleRegion,
					Scale:       40.0, // TODO Make this configurable by reading the editor engine state.
					Orientation: 0,
					Position:    n.getEffectivePosition()}})
	} else {
		e.Reset()
		e.Lines = []HypotheticalLine{
			HypotheticalLine{
				Color: mgl32.Vec3{1.0, 1.0, 0.0},
				Line: [2]mgl32.Vec2{
					n.getEffectivePosition(),
					n.roadLineState.firstNode}}}
	}
}

func (e *HypotheticalActions) computeDrawIndicator(n *Engine) {
	e.setSingleRegion(
		HypotheticalRegion{
			Color: mgl32.Vec3{0.0, 1.0, 1.0},
			Region: commonMath.Region{
				RegionType:  commonMath.CircleRegion,
				Scale:       30.0, // TODO Make this configurable by reading the editor engine state.
				Orientation: 0,
				Position:    n.lastBoardPos}})
}

// Updates the hypotheticals to be applicable to the current edit mode.
func (e *HypotheticalActions) ComputeHypotheticalRegion(n *Engine, engineState *editorEngine.State) {
	if engineState.Mode == editorEngine.Add {
		if engineState.InAddMode == editorEngine.PowerPlant {
			e.computePowerPlantHypotheticalRegion(n)
		} else if engineState.InAddMode == editorEngine.PowerLine {
			e.computePowerLineHypotheticalRegion(n)
		} else if engineState.InAddMode == editorEngine.RoadLine {
			e.computeRoadLineHypotheticalRegion(n)
		}
	} else if engineState.Mode == editorEngine.Draw {
		e.computeDrawIndicator(n)
	} else {
		e.Reset()
	}
}
