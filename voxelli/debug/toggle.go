package debug

import (
	"go-experiments/voxelli/input"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var wasDebugPressed bool = false
var isDebug bool = false

func CheckDebugToggle() {
	if !wasDebugPressed && input.PressedKeys[glfw.KeyT] {
		wasDebugPressed = true
		isDebug = !isDebug
	}

	if wasDebugPressed && !input.PressedKeys[glfw.KeyT] {
		wasDebugPressed = false
	}
}

func IsDebug() bool {
	return isDebug
}
