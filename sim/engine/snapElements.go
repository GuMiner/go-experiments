package engine

import (
	"go-experiments/common/commonmath"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/element"
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/mathgl/mgl32"
)

type SnapElements struct {
	snappedNode    *element.ElementWithDistance
	snappedLine    *element.LineElementWithDistance
	snappedGridPos *mgl32.Vec2
}

func NewSnapElements() SnapElements {
	s := SnapElements{}
	s.Reset()
	return s
}

func (s *SnapElements) ComputeSnappedSnapElements(boardPos mgl32.Vec2, elementFinder *element.ElementFinder) {
	s.snappedNode = nil
	if editorEngine.EngineState.SnapToElements &&
		editorEngine.EngineState.Mode == editorEngine.Add &&
		editorEngine.EngineState.InAddMode == editorEngine.PowerLine {

		elements := elementFinder.KNearest(boardPos, config.Config.Draw.SnapNodeCount)
		for _, elem := range elements {
			if elem.Distance < config.Config.Draw.MinSnapNodeDistance {
				s.snappedNode = &elem
				break
			}
		}
	}

	s.snappedGridPos = nil
	if editorEngine.EngineState.SnapToGrid {
		snapGridResolution := float32(config.Config.Snap.SnapGridResolution)
		offsetBoardPos := boardPos.Add(mgl32.Vec2{snapGridResolution / 2, snapGridResolution / 2})
		snappedIntPosition := commonMath.IntVec2{int(offsetBoardPos.X() / snapGridResolution), int(offsetBoardPos.Y() / snapGridResolution)}
		snappedPos := mgl32.Vec2{float32(snappedIntPosition.X()), float32(snappedIntPosition.Y())}.Mul(snapGridResolution)
		s.snappedGridPos = &snappedPos
	}
}

func (s *SnapElements) GetSnappedNode() *element.ElementWithDistance {
	return s.snappedNode
}

func (s *SnapElements) GetSnappedLine() *element.LineElementWithDistance {
	return s.snappedLine
}

func (s *SnapElements) Reset() {
	s.snappedLine = nil
	s.snappedLine = nil
}
