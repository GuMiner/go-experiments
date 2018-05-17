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

	SnapToGridKey
	SnapToAngleKey
	SnapToElementsKey

	SelectModeKey
	AddModeKey
	DrawModeKey

	PowerPlantAddModeKey
	PowerLineAddModeKey

	CoalPlantKey
	NuclearPlantKey
	NaturalGasPlantKey
	WindPlantKey
	SolarPlantKey
	GeothermalPlantKey

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
	keyMap[CoalPlantKey] = glfw.Key1
	keyMap[NuclearPlantKey] = glfw.Key2
	keyMap[NaturalGasPlantKey] = glfw.Key3
	keyMap[WindPlantKey] = glfw.Key4
	keyMap[SolarPlantKey] = glfw.Key5
	keyMap[GeothermalPlantKey] = glfw.Key6

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

	keyMap[SnapToGridKey] = glfw.Key8
	keyMap[SnapToAngleKey] = glfw.Key9
	keyMap[SnapToElementsKey] = glfw.Key0

	keyMap[SelectModeKey] = glfw.KeyS
	keyMap[AddModeKey] = glfw.KeyA
	keyMap[DrawModeKey] = glfw.KeyD

	keyMap[PowerPlantAddModeKey] = glfw.KeyP
	keyMap[PowerLineAddModeKey] = glfw.KeyL

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
