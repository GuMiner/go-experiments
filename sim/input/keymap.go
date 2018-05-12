package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Defines the key assignments of each key
type KeyAssignment int

// Ordering is important! Reordering these changes their numerical values
const (
	MoveLeft KeyAssignment = iota
	MoveRight
	MoveUp
	MoveDown

	SnapToGrid
	SnapToAngle
	SnapToElements

	SelectMode
	AddMode

	PowerPlantAddMode
	PowerLineAddMode
)

const keyMapCacheName = "keymap"

var keyMap map[KeyAssignment]glfw.Key = make(map[KeyAssignment]glfw.Key)

func CreateDefaultKeyMap() {
	keyMap[MoveLeft] = glfw.KeyLeft
	keyMap[MoveRight] = glfw.KeyRight
	keyMap[MoveUp] = glfw.KeyUp
	keyMap[MoveDown] = glfw.KeyDown

	keyMap[SnapToGrid] = glfw.Key8
	keyMap[SnapToAngle] = glfw.Key9
	keyMap[SnapToElements] = glfw.Key0

	keyMap[SelectMode] = glfw.KeyS
	keyMap[AddMode] = glfw.KeyA

	keyMap[PowerPlantAddMode] = glfw.KeyP
	keyMap[PowerLineAddMode] = glfw.KeyL
}

func IsPressed(key KeyAssignment) bool {
	return pressedKeys[keyMap[key]]
}

func IsTyped(key KeyAssignment) bool {
	isTyped := typedKeys[keyMap[key]]
	if isTyped {
		typedKeys[keyMap[key]] = false
	}

	return isTyped
}
