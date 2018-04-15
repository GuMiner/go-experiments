package input

import (
	"go-experiments/voxelli/cache"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Defines the key assignments of each key
type KeyAssignment int

// Ordering is important! Reordering these changes their numerical values
const (
	LookLeft KeyAssignment = iota
	LookRight
	LookUp
	LookDown

	RotateClockwise
	RotateCounterClockwise

	MoveForwards
	MoveBackwards
	MoveLeft
	MoveRight
	MoveUp
	MoveDown

	ToggleDebug
	ToggleWireframe
	ToggleColorOverlay
)

const keyMapCacheName = "keymap"

var keyMap map[KeyAssignment]glfw.Key = make(map[KeyAssignment]glfw.Key)

func createDefaultKeyMap() {
	keyMap[LookLeft] = glfw.KeyLeft
	keyMap[LookRight] = glfw.KeyRight
	keyMap[LookUp] = glfw.KeyUp
	keyMap[LookDown] = glfw.KeyDown

	keyMap[RotateClockwise] = glfw.KeyE
	keyMap[RotateCounterClockwise] = glfw.KeyD

	keyMap[MoveForwards] = glfw.KeyA
	keyMap[MoveBackwards] = glfw.KeyZ
	keyMap[MoveLeft] = glfw.KeyQ
	keyMap[MoveRight] = glfw.KeyW
	keyMap[MoveUp] = glfw.KeyS
	keyMap[MoveDown] = glfw.KeyX

	keyMap[ToggleDebug] = glfw.KeyT
	keyMap[ToggleWireframe] = glfw.KeyR
	keyMap[ToggleColorOverlay] = glfw.KeyC
}

func LoadKeyAssignments() {
	cacheMiss := cache.LoadFromCache(keyMapCacheName, false, &keyMap)
	if cacheMiss {
		createDefaultKeyMap()
	}
}

func SaveKeyAssignments() {
	cache.SaveToCache(keyMapCacheName, keyMap)
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
