package ui

import (
	"go-experiments/sim/input/editorEngine"
	"go-experiments/sim/ui/overlay"
	"go-experiments/sim/ui/region"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type UiInfrastructure struct {
	OverlayProgram *overlay.OverlayShaderProgram
	RegionProgram  *region.RegionShaderProgram
}

var Ui UiInfrastructure

// Defines common UI initialization, for both 2D and 3D rendering modes.
func Init(window *glfw.Window) {
	Ui.OverlayProgram = overlay.NewOverlayShaderProgram()
	Ui.RegionProgram = region.NewRegionShaderProgram()

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
	Ui.OverlayProgram.Delete()
	Ui.RegionProgram.Delete()

	destroyCustomCursors()
}
