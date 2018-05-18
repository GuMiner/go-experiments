package commonOpenGl

import (
	"go-experiments/common/commonconfig"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var windowSize mgl32.Vec2

func InitViewport() {
	windowSize = mgl32.Vec2{
		float32(commonConfig.Config.Window.Width),
		float32(commonConfig.Config.Window.Height)}
}

func ResizeViewport(window *glfw.Window, width int, height int) {
	windowSize = mgl32.Vec2{float32(width), float32(height)}
}

func GetWindowSize() mgl32.Vec2 {
	return windowSize
}

func ResetViewport() {
	gl.Viewport(0, 0, int32(windowSize.X()), int32(windowSize.Y()))
}
