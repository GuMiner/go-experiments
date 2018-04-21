package viewport

import (
	"go-experiments/voxelli/config"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const aspectRatio float32 = 1.0

var windowSize mgl32.Vec2

func Init() {
	windowSize = mgl32.Vec2{
		float32(config.Config.Window.Width),
		float32(config.Config.Window.Height)}
}

func HandleResize(window *glfw.Window, width int, height int) {
	windowSize = mgl32.Vec2{float32(width), float32(height)}
}

func GetWidth() float32 {
	return windowSize.X()
}

func GetHeight() float32 {
	return windowSize.Y()
}

func Reset() {
	gl.Viewport(0, 0, int32(windowSize.X()), int32(windowSize.Y()))
}
