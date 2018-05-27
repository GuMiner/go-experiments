package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type KeyAssignment int

const (
	MoveLeftKey KeyAssignment = iota
	MoveRightKey
	MoveUpKey
	MoveDownKey

	PauseKey
	CancelKey

	SnapToGridKey
	SnapToAngleKey
	SnapToElementsKey

	SelectModeKey
	AddModeKey
	DrawModeKey

	PowerPlantAddModeKey
	PowerLineAddModeKey
	RoadLineAddModeKey

	ItemAdd1Key
	ItemAdd2Key
	ItemAdd3Key
	ItemAdd4Key
	ItemAdd5Key
	ItemAdd6Key

	TerrainFlattenKey
	TerrainSharpenKey
	TerrainTreesKey
	TerrainShrubsKey
	TerrainHillsKey
	TerrainValleysKey
)

const keyMapCacheName = "keymap"

var keyMap map[KeyAssignment]glfw.Key = make(map[KeyAssignment]glfw.Key)

// Sub options are 1-7
func createSubOptionsKeyMap() {
	keyMap[ItemAdd1Key] = glfw.Key1
	keyMap[ItemAdd2Key] = glfw.Key2
	keyMap[ItemAdd3Key] = glfw.Key3
	keyMap[ItemAdd4Key] = glfw.Key4
	keyMap[ItemAdd5Key] = glfw.Key5
	keyMap[ItemAdd6Key] = glfw.Key6

	keyMap[TerrainFlattenKey] = glfw.Key1
	keyMap[TerrainSharpenKey] = glfw.Key2
	keyMap[TerrainTreesKey] = glfw.Key3
	keyMap[TerrainShrubsKey] = glfw.Key4
	keyMap[TerrainHillsKey] = glfw.Key5
	keyMap[TerrainValleysKey] = glfw.Key6
}

func CreateDefaultKeyMap() {
	keyMap[MoveLeftKey] = glfw.KeyLeft
	keyMap[MoveRightKey] = glfw.KeyRight
	keyMap[MoveUpKey] = glfw.KeyUp
	keyMap[MoveDownKey] = glfw.KeyDown

	keyMap[PauseKey] = glfw.KeySpace
	keyMap[CancelKey] = glfw.KeyEscape

	keyMap[SnapToGridKey] = glfw.Key8
	keyMap[SnapToAngleKey] = glfw.Key9
	keyMap[SnapToElementsKey] = glfw.Key0

	keyMap[SelectModeKey] = glfw.KeyS
	keyMap[AddModeKey] = glfw.KeyA
	keyMap[DrawModeKey] = glfw.KeyD

	keyMap[PowerPlantAddModeKey] = glfw.KeyP
	keyMap[PowerLineAddModeKey] = glfw.KeyL
	keyMap[RoadLineAddModeKey] = glfw.KeyR

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
