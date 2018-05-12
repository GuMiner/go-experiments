package ui

import (
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Defines common UI initialization, for both 2D and 3D rendering modes.
func Init(window *glfw.Window) {
	initCustomCursors(window)
}

func UpdateEditorState(engineState editorEngine.State, window *glfw.Window) {
	cursor := Selection

	if engineState.Mode == editorEngine.Select {
		cursor = Selection
	} else { // Add mode for now
		if engineState.InAddMode == editorEngine.PowerLine {
			cursor = PowerLineAdd
		} else if engineState.InAddMode == editorEngine.PowerPlant {
			cursor = PowerPlantAdd
		}
	}

	setCursor(cursor, window)
}

func Delete() {
	destroyCustomCursors()
}
