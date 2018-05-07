package ui

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

var customCursor *glfw.Cursor

func initCustomCursors(window *glfw.Window) {
	// TODO: Custom cursors for actions. Flesh this more out once we actually start applying actions.
	customCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	window.SetCursor(customCursor)
}

func destroyCustomCursors() {
	customCursor.Destroy()
}
