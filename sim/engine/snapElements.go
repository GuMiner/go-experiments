package engine

import "go-experiments/sim/engine/element"

type SnapElements struct {
	snappedNode *element.ElementWithDistance
	snappedLine *element.LineElementWithDistance
}

func NewSnapElements() SnapElements {
	s := SnapElements{}
	s.Reset()
	return s
}

func (s *SnapElements) ComputeSnappedSnapElements() {

}

func (s *SnapElements) Reset() {
	s.snappedLine = nil
	s.snappedLine = nil
}
