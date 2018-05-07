package editorEngine

import (
	"fmt"
	"go-experiments/sim/input"
)

// TODO: Modes should really be saved in game state on a per-game basis.
var EngineState State

func Init() {
	EngineState = State{
		Mode:           Select,
		SnapToGrid:     true,
		SnapToElements: true,
		SnapToAngle:    false}
}

// Update toggles and the current edit mode.
func Update() {
	if input.IsTyped(input.SnapToGrid) {
		EngineState.SnapToGrid = !EngineState.SnapToGrid
		fmt.Printf("Snap to grid: %v\n", EngineState.SnapToGrid)
	}

	if input.IsTyped(input.SnapToAngle) {
		EngineState.SnapToAngle = !EngineState.SnapToAngle
		fmt.Printf("Snap to angle: %v\n", EngineState.SnapToAngle)
	}

	if input.IsTyped(input.SnapToElements) {
		EngineState.SnapToElements = !EngineState.SnapToElements
		fmt.Printf("Snap to elements: %v\n", EngineState.SnapToElements)
	}

	if input.IsTyped(input.SelectMode) {
		EngineState.Mode = Select
		fmt.Printf("Entered selection mode\n")
	}

	if input.IsTyped(input.AddMode) {
		EngineState.Mode = Add
		fmt.Printf("Entered addition mode\n")
	}
}
