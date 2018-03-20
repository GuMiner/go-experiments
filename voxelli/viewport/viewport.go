package viewport

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const aspectRatio float32 = 1.0

var windowSize mgl32.Vec2 = mgl32.Vec2{800, 600}
var perspectiveMatrixUpdated bool = false

func HandleResize(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	windowSize = mgl32.Vec2{float32(width), float32(height)}
	perspectiveMatrixUpdated = false
}

func GetWidth() float32 {
	return windowSize.X()
}

func GetHeight() float32 {
	return windowSize.Y()
}

func PerspectiveMatrixUpdated() bool {
	wasUpdated := perspectiveMatrixUpdated
	perspectiveMatrixUpdated = true
	return wasUpdated
}
