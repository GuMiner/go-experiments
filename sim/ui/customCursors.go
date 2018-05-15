package ui

import (
	"go-experiments/common/io"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type CustomCursorType int

const (
	Selection CustomCursorType = iota
	PowerPlantAdd
	PowerLineAdd
)

type CustomCursors struct {
	cursors map[CustomCursorType]*glfw.Cursor
}

var customCursors CustomCursors

func initCustomCursors(window *glfw.Window) {
	customCursors = CustomCursors{
		cursors: make(map[CustomCursorType]*glfw.Cursor)}

	customCursors.cursors[Selection] = glfw.CreateStandardCursor(glfw.ArrowCursor)

	powerPlantImage := commonIo.ReadImageFromFile("data/cursors/PowerPlant.png")
	powerLineImage := commonIo.ReadImageFromFile("data/cursors/PowerLine.png")
	customCursors.cursors[PowerPlantAdd] = glfw.CreateCursor(powerPlantImage, 0, 0)
	customCursors.cursors[PowerLineAdd] = glfw.CreateCursor(powerLineImage, 0, 0)

	window.SetCursor(customCursors.cursors[Selection])
}

func setCursor(cursorType CustomCursorType, window *glfw.Window) {
	window.SetCursor(customCursors.cursors[cursorType])
}

func destroyCustomCursors() {
	for _, cursor := range customCursors.cursors {
		cursor.Destroy()
	}
}
