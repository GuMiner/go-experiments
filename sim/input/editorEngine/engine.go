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
		InAddMode:      PowerPlant,
		SnapToGrid:     true,
		SnapToElements: true,
		SnapToAngle:    false}
}

func updateEngineState() bool {
	updated := false
	if input.IsTyped(input.SnapToGrid) {
		EngineState.SnapToGrid = !EngineState.SnapToGrid
		fmt.Printf("Snap to grid: %v\n", EngineState.SnapToGrid)
		updated = true
	}

	if input.IsTyped(input.SnapToAngle) {
		EngineState.SnapToAngle = !EngineState.SnapToAngle
		fmt.Printf("Snap to angle: %v\n", EngineState.SnapToAngle)
		updated = true
	}

	if input.IsTyped(input.SnapToElements) {
		EngineState.SnapToElements = !EngineState.SnapToElements
		fmt.Printf("Snap to elements: %v\n", EngineState.SnapToElements)
		updated = true
	}

	if input.IsTyped(input.SelectMode) {
		EngineState.Mode = Select
		fmt.Printf("Entered selection mode\n")
		updated = true
	}

	if input.IsTyped(input.AddMode) {
		EngineState.Mode = Add
		fmt.Printf("Entered addition mode\n")
		updated = true
	}

	return updated
}

func updateEngineAddState() bool {
	updated := false
	if input.IsTyped(input.PowerPlantAddMode) {
		EngineState.InAddMode = PowerPlant
		fmt.Printf("Entered power plant addition mode (applies when in Add mode)\n")
		updated = true
	}

	if input.IsTyped(input.PowerLineAddMode) {
		EngineState.InAddMode = PowerLine
		fmt.Printf("Entered power line addition mode (applies when in Add mode)\n")
		updated = true
	}

	return updated
}

// Update toggles and the current edit mode.
func Update() bool {
	updatedEngine := updateEngineState()
	updatedAdd := updateEngineAddState()
	return updatedEngine || updatedAdd
}
