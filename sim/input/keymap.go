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

	CoalPlantKey
	NuclearPlantKey
	NaturalGasPlantKey
	WindPlantKey
	SolarPlantKey
	GeothermalPlantKey
)

const keyMapCacheName = "keymap"

var keyMap map[KeyAssignment]glfw.Key = make(map[KeyAssignment]glfw.Key)

// Sub options are 1-7
func createSubOptionsKeyMap() {
	keyMap[CoalPlantKey] = glfw.Key1
	keyMap[NuclearPlantKey] = glfw.Key2
	keyMap[NaturalGasPlantKey] = glfw.Key3
	keyMap[WindPlantKey] = glfw.Key4
	keyMap[SolarPlantKey] = glfw.Key5
	keyMap[GeothermalPlantKey] = glfw.Key6
}

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

	createSubOptionsKeyMap()
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
