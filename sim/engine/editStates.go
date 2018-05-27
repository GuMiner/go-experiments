package engine

import (
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/mathgl/mgl32"
)

type PowerLineEditState struct {
	wasInEditState bool

	hasFirstNode     bool
	firstNode        mgl32.Vec2
	firstNodeElement int
}

func NewPowerLineEditState() *PowerLineEditState {
	s := PowerLineEditState{}
	s.Reset()
	return &s
}

// TODO: Deduplicate
type RoadLineEditState struct {
	wasInEditState bool

	hasFirstNode     bool
	firstNode        mgl32.Vec2
	firstNodeElement int
}

func NewRoadLineEditState() *RoadLineEditState {
	s := RoadLineEditState{}
	s.Reset()
	return &s
}

func (p *PowerLineEditState) InPowerLineState(engineState *editorEngine.State) bool {
	return engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerLine
}

func (p *RoadLineEditState) InRoadLineState(engineState *editorEngine.State) bool {
	return engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.RoadLine
}

func (p *PowerLineEditState) EnterOrExitEditMode(engineState *editorEngine.State) {
	inPowerLineEdit := p.InPowerLineState(engineState)
	if inPowerLineEdit && !p.wasInEditState {
		p.Reset()
		p.wasInEditState = true
	} else if !inPowerLineEdit && p.wasInEditState {
		p.wasInEditState = false
	}
}

func (p *RoadLineEditState) EnterOrExitEditMode(engineState *editorEngine.State) {
	inRoadLineState := p.InRoadLineState(engineState)
	if inRoadLineState && !p.wasInEditState {
		p.Reset()
		p.wasInEditState = true
	} else if !inRoadLineState && p.wasInEditState {
		p.wasInEditState = false
	}
}

func (p *PowerLineEditState) Reset() {
	p.wasInEditState = false
	p.hasFirstNode = false
}

func (p *RoadLineEditState) Reset() {
	p.wasInEditState = false
	p.hasFirstNode = false
}
