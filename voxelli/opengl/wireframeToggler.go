package opengl

import (
	"go-experiments/voxelli/input"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var isWireframe bool = false

func CheckWireframeToggle() {
	if input.IsTyped(input.ToggleWireframe) {
		isWireframe = !isWireframe
		if isWireframe {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}
	}
}
