package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var anyEvent bool = false
var MousePos mgl32.Vec2 = mgl32.Vec2{0, 0}
var PressedKeys map[glfw.Key]bool = make(map[glfw.Key]bool)
var PressedButtons map[glfw.MouseButton]bool = make(map[glfw.MouseButton]bool)

func HandleMouseMove(window *glfw.Window, xPos float64, yPos float64) {
	MousePos = mgl32.Vec2{float32(xPos), float32(yPos)}
	anyEvent = true
}

func HandleMouseButton(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		PressedButtons[button] = true
		anyEvent = true
	case glfw.Release:
		PressedButtons[button] = false
		anyEvent = true
	}
}

func HandleKeyInput(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		PressedKeys[key] = true
		anyEvent = true
	case glfw.Release:
		PressedKeys[key] = false
		anyEvent = true
	}
}

func AnyEvent() bool {
	if anyEvent {
		anyEvent = false
		return true
	}

	return false
}
