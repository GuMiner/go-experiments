package ui

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Defines common UI initialization, for both 2D and 3D rendering modes.
func Init(window *glfw.Window) {
	initCustomCursors(window)
}

func Delete() {
	destroyCustomCursors()
}
