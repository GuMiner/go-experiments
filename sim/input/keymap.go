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
)

const keyMapCacheName = "keymap"

var keyMap map[KeyAssignment]glfw.Key = make(map[KeyAssignment]glfw.Key)

func CreateDefaultKeyMap() {
	keyMap[MoveLeft] = glfw.KeyA
	keyMap[MoveRight] = glfw.KeyD
	keyMap[MoveUp] = glfw.KeyW
	keyMap[MoveDown] = glfw.KeyS
}

// func LoadKeyAssignments() {
// 	cacheMiss := cache.LoadFromCache(keyMapCacheName, false, &keyMap)
// 	if cacheMiss {
// 		createDefaultKeyMap()
// 	}
// }
//
// func SaveKeyAssignments() {
// 	cache.SaveToCache(keyMapCacheName, keyMap)
// }

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
