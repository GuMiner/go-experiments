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

func (p *PowerLineEditState) InPowerLineState(engineState *editorEngine.State) bool {
	return engineState.Mode == editorEngine.Add && engineState.InAddMode == editorEngine.PowerLine
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

func (p *PowerLineEditState) Reset() {
	p.wasInEditState = false
	p.hasFirstNode = false
}
