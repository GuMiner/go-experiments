package commonOpenGl

import (
	"go-experiments/common/config"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const aspectRatio float32 = 1.0

var windowSize mgl32.Vec2

func InitViewport() {
	windowSize = mgl32.Vec2{
		float32(commonConfig.Config.Window.Width),
		float32(commonConfig.Config.Window.Height)}
}

func ResizeViewport(window *glfw.Window, width int, height int) {
	windowSize = mgl32.Vec2{float32(width), float32(height)}
}

func GetViewportWidth() float32 {
	return windowSize.X()
}

func GetViewportHeight() float32 {
	return windowSize.Y()
}

func ResetViewport() {
	gl.Viewport(0, 0, int32(windowSize.X()), int32(windowSize.Y()))
}