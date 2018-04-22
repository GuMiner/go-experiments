package main

import (
	"go-experiments/common/color"
	"go-experiments/common/config"
	"go-experiments/common/diagnostics"
	"go-experiments/common/opengl"
	"go-experiments/common/shadow"
	"go-experiments/voxelli/input"
	"go-experiments/voxelli/renderer"

	"github.com/go-gl/mathgl/mgl32"

	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

func setInputCallbacks(window *glfw.Window) {
	window.SetFramebufferSizeCallback(commonOpenGl.ResizeViewport)
	window.SetCursorPosCallback(input.HandleMouseMove)
	window.SetMouseButtonCallback(input.HandleMouseButton)
	window.SetKeyCallback(input.HandleKeyInput)
}

func main() {
	commonConfig.Load("./data/commonConfig.json")

	commonOpenGl.InitGlfw()
	defer glfw.Terminate()

	commonOpenGl.InitViewport()
	window, err := glfw.CreateWindow(
		int(commonOpenGl.GetViewportWidth()),
		int(commonOpenGl.GetViewportHeight()),
		commonConfig.Config.Window.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	setInputCallbacks(window)
	commonOpenGl.ConfigureOpenGl()

	//input.LoadKeyAssignments()
	//defer input.SaveKeyAssignments()

	// Create renderers
	shadowBuffer := shadow.NewShadowBuffer()
	defer shadowBuffer.Delete()

	commonDiagnostics.InitCube()
	defer commonDiagnostics.DeleteCube()

	commonColor.InitializeColorGradient(
		commonConfig.Config.ColorGradient.Steps,
		commonConfig.Config.ColorGradient.Saturation,
		commonConfig.Config.ColorGradient.Luminosity)

	startTime := time.Now()
	frameTime := float32(0.1)
	lastElapsed := float32(0.0)
	elapsed := lastElapsed

	// Update
	update := func() {
		lastElapsed = elapsed
		elapsed = float32(time.Since(startTime)) / float32(time.Second)
		frameTime = elapsed - lastElapsed

		// fpsCounter.Update(frameTime)
		// opengl.CheckWireframeToggle()
		// diagnostics.CheckDebugToggle()

		// camera.Update(frameTime)

		glfw.PollEvents()
	}

	render := func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		model := mgl32.Translate3D(0, 0, 0)
		commonDiagnostics.Render(
			mgl32.Vec4{0, 1, 0, 1},
			&model)
	}

	for !window.ShouldClose() {
		update()

		// Render to shadow buffer
		// projection, cameraMatrix := shadowBuffer.PrepareCamera()

		// TODO: Update projection and camera of renderable items

		shadowBuffer.RenderToBuffer(func() {
			gl.Clear(gl.DEPTH_BUFFER_BIT)
			gl.CullFace(gl.FRONT)

			renderer.EnableDepthModeOnly()
			// TODO: Render depth only
			renderer.DisableDepthModeOnly()

			gl.CullFace(gl.BACK)
		})

		// Prepare for normal rendering...
		/*shadowBiasMatrix := mgl32.Mat4FromRows(
		mgl32.Vec4{0.5, 0, 0, 0.5},
		mgl32.Vec4{0, 0.5, 0, 0.5},
		mgl32.Vec4{0, 0, 0.5, 0.5},
		mgl32.Vec4{0, 0, 0, 1.0})*/

		// partialShadowMatrix := shadowBiasMatrix.Mul4(projection.Mul4(cameraMatrix))
		// TODO: Pass partial shadow matrix and shadow image into program.

		// Render the full display.
		commonOpenGl.ResetViewport()

		projection := mgl32.Perspective(
			mgl32.DegToRad(commonConfig.Config.Perspective.FovY),
			commonOpenGl.GetViewportWidth()/commonOpenGl.GetViewportHeight(),
			commonConfig.Config.Perspective.Near,
			commonConfig.Config.Perspective.Far)
		cube := commonDiagnostics.GetCube()
		cube.UpdateProjection(&projection)

		cameraMatrix := mgl32.LookAtV(
			mgl32.Vec3{1, 1, 1},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 1})
		cube.UpdateCamera(&cameraMatrix)

		render()
		window.SwapBuffers()
	}
}
