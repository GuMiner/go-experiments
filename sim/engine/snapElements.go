package engine

import (
	"go-experiments/sim/config"
	"go-experiments/sim/engine/element"
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/mathgl/mgl32"
)

type SnapElements struct {
	snappedNode *element.ElementWithDistance
	snappedLine *element.LineElementWithDistance
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
