package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func initGlfw() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
}

var windowSize mgl32.Vec2 = mgl32.Vec2{800, 600}

func handleResize(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	windowSize = mgl32.Vec2{float32(width), float32(height)}
}

var mousePos mgl32.Vec2 = mgl32.Vec2{0, 0}

func handleMouseMove(window *glfw.Window, xPos float64, yPos float64) {
	mousePos = mgl32.Vec2{float32(xPos), float32(yPos)}
}

func setInputCallbacks(window *glfw.Window) {
	window.SetFramebufferSizeCallback(handleResize)
	window.SetCursorPosCallback(handleMouseMove)
}

func main() {
	initGlfw()
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(int(windowSize.X()), int(windowSize.Y()), "Fractal", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	setInputCallbacks(window)
	configureOpenGl()

	program := createProgram("./shaders/juliaFractal")
	defer gl.DeleteProgram(program)

	// Get locations of everything used in this program.
	mouseUniformLoc := gl.GetUniformLocation(program, gl.Str("c\x00"))
	timeLoc := gl.GetUniformLocation(program, gl.Str("time\x00"))
	fractalGradientLoc := gl.GetUniformLocation(program, gl.Str("fractalGradient\x00"))
	maxIterationsLoc := gl.GetUniformLocation(program, gl.Str("maxIterations\x00"))

	var maxIterations int32 = 1000
	rainbow := NewRainbowTexture(maxIterations)
	defer rainbow.Delete()

	quad := NewFullScreenQuad()
	defer quad.Delete()

	// Enable our program and do first-time uniform setup
	gl.UseProgram(program)

	gl.ActiveTexture(rainbow.openGlTextureSlot)
	gl.BindTexture(gl.TEXTURE_1D, rainbow.fractalGradientTextureID)
	gl.Uniform1i(fractalGradientLoc, int32(rainbow.openGlTextureSlot)-gl.TEXTURE0) // Yes it is zero.

	gl.Uniform1i(maxIterationsLoc, maxIterations)

	startTime := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update uniforms that need updating
		elapsed := time.Since(startTime)
		gl.Uniform1f(timeLoc, float32(elapsed)/float32(time.Second))

		normalizedPos := mgl32.Vec2{-0.5, -0.5}.Add(mgl32.Vec2{mousePos.X() / windowSize.X(), mousePos.Y() / windowSize.Y()})

		gl.Uniform2f(mouseUniformLoc, normalizedPos.X(), normalizedPos.Y())

		quad.Render()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
