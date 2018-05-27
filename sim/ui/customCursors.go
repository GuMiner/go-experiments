package ui

import (
	"go-experiments/common/commonio"
	"go-experiments/sim/input/editorEngine"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type CustomCursorType int

const (
	Selection CustomCursorType = iota
	PowerPlantAdd
	PowerLineAdd
	RoadLineAdd

	TerrainFlatten
	TerrainSharpen
	TerrainTrees
	TerrainShrubs
	TerrainHills
	TerrainValleys
)

type CustomCursors struct {
	cursors map[CustomCursorType]*glfw.Cursor
}

var drawModeCursors map[editorEngine.EditorDrawMode]CustomCursorType = make(map[editorEngine.EditorDrawMode]CustomCursorType)

var customCursors CustomCursors

func initDrawModeCursors(customCursors *CustomCursors) {
	type DrawModeCursor struct {
		location     string
		cursorType   CustomCursorType
		drawModeType editorEngine.EditorDrawMode
	}

	cursorsToLoad := []DrawModeCursor{
		DrawModeCursor{"data/cursors/draw/TerrainFlatten.png", TerrainFlatten, editorEngine.TerrainFlatten},
		DrawModeCursor{"data/cursors/draw/TerrainSharpen.png", TerrainSharpen, editorEngine.TerrainSharpen},
		DrawModeCursor{"data/cursors/draw/TerrainTrees.png", TerrainTrees, editorEngine.TerrainTrees},
		DrawModeCursor{"data/cursors/draw/TerrainShrubs.png", TerrainShrubs, editorEngine.TerrainShrubs},
		DrawModeCursor{"data/cursors/draw/TerrainHills.png", TerrainHills, editorEngine.TerrainHills},
		DrawModeCursor{"data/cursors/draw/TerrainValleys.png", TerrainValleys, editorEngine.TerrainValleys}}

	for _, cursor := range cursorsToLoad {
		cursorImage := commonIo.ReadImageFromFile(cursor.location)
		customCursors.cursors[cursor.cursorType] = glfw.CreateCursor(cursorImage, 0, 0)
		drawModeCursors[cursor.drawModeType] = cursor.cursorType
	}
}

func initCustomCursors(window *glfw.Window) {
	customCursors = CustomCursors{
		cursors: make(map[CustomCursorType]*glfw.Cursor)}

	customCursors.cursors[Selection] = glfw.CreateStandardCursor(glfw.ArrowCursor)

	powerPlantImage := commonIo.ReadImageFromFile("data/cursors/PowerPlant.png")
	powerLineImage := commonIo.ReadImageFromFile("data/cursors/PowerLine.png")
	roadLineImage := commonIo.ReadImageFromFile("data/cursors/RoadLine.png")
	customCursors.cursors[PowerPlantAdd] = glfw.CreateCursor(powerPlantImage, 0, 0)
	customCursors.cursors[PowerLineAdd] = glfw.CreateCursor(powerLineImage, 0, 0)
	customCursors.cursors[RoadLineAdd] = glfw.CreateCursor(roadLineImage, 0, 0)
	initDrawModeCursors(&customCursors)

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
