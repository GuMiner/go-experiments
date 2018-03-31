package opengl

import (
	"go-experiments/voxelli/input"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var wasTogglePressed bool = false
var isWireframe bool = false

func CheckWireframeToggle() {
	if !wasTogglePressed && input.PressedKeys[glfw.KeyR] {
		wasTogglePressed = true
		isWireframe = !isWireframe
		if isWireframe {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}
	}

	if wasTogglePressed && !input.PressedKeys[glfw.KeyR] {
		wasTogglePressed = false
	}
}
