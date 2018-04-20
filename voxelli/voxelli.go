package main

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.

import (
	"fmt"
	"go-experiments/voxelli/color"
	"go-experiments/voxelli/diagnostics"
	"go-experiments/voxelli/input"
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/text"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/viewport"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

func setInputCallbacks(window *glfw.Window) {
	window.SetFramebufferSizeCallback(viewport.HandleResize)
	window.SetCursorPosCallback(input.HandleMouseMove)
	window.SetMouseButtonCallback(input.HandleMouseButton)
	window.SetKeyCallback(input.HandleKeyInput)
}

func main() {
	opengl.InitGlfw()
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(int(viewport.GetWidth()), int(viewport.GetHeight()), "Voxelli", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	setInputCallbacks(window)
	opengl.ConfigureOpenGl()

	input.LoadKeyAssignments()
	defer input.SaveKeyAssignments()

	// Create renderers
	shadowBuffer := NewShadowBuffer()
	defer shadowBuffer.Delete()

	diagnostics.InitCube()
	defer diagnostics.DeleteCube()

	color.InitializeColorGradient(400, 1.0, 0.5)

	voxelArrayObjectRenderer := renderer.NewVoxelArrayObjectRenderer()
	defer voxelArrayObjectRenderer.Delete()

	textRenderer := text.NewTextRenderer("./data/font/DejaVuSans.ttf")
	defer textRenderer.Delete()

	testSentence := text.NewSentence(textRenderer, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 0})
	fpsSentence := text.NewSentence(textRenderer, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	var renderers []renderer.Renderer
	renderers = append(renderers, voxelArrayObjectRenderer)
	renderers = append(renderers, textRenderer)
	renderers = append(renderers, diagnostics.GetCube())

	camera := NewCamera(mgl32.Vec3{140, 300, 300}, mgl32.Vec3{-1, 0, 0}, mgl32.Vec3{0, 0, 1})
	defer camera.CachePosition()

	InitSimulation(voxelArrayObjectRenderer)
	defer DeleteSimulation()

	startTime := time.Now()
	frameTime := float32(0.1)
	lastElapsed := float32(0.0)
	elapsed := lastElapsed

	// Update
	update := func() {
		lastElapsed = elapsed
		elapsed = float32(time.Since(startTime)) / float32(time.Second)
		frameTime = elapsed - lastElapsed

		opengl.CheckWireframeToggle()
		diagnostics.CheckDebugToggle()
		vehicle.CheckColorOverlayToggle()

		camera.Update(frameTime)

		UpdateSimulation(frameTime, elapsed)
		glfw.PollEvents()
	}

	render := func() {
		RenderSimulation(voxelArrayObjectRenderer)

		textModelMatrix := mgl32.Translate3D(20, 20, 20).Mul4(mgl32.Scale3D(3, 3, 1))
		testSentence.Render("Hello world!.,J12359?/~╠>B☢>124$", &textModelMatrix, true)

		fpsModelMatrix := mgl32.Translate3D(0, 0, 20).Mul4(mgl32.Scale3D(3, 3, 1))
		fpsSentence.Render(fmt.Sprintf("FPS: %.2f", 1.0/frameTime), &fpsModelMatrix, true)
	}

	for !window.ShouldClose() {
		update()

		// Render to shadow buffer
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.Viewport(0, 0, shadowBuffer.Width, shadowBuffer.Height)

		// Hardcoded for the following properties:
		// 1. Good shadow angle
		// 2. Good usage of the depth buffer range
		projection := mgl32.Ortho(-120, 120, -120, 120, 760, 1000)
		renderer.UpdateProjections(renderers, &projection)

		position := mgl32.Vec3{-5.7113113, -642.92566, 476.05392}
		cameraMatrix := mgl32.LookAtV(
			position,
			position.Add(mgl32.Vec3{0.117314756, 0.8421086, -0.52639383}),
			mgl32.Vec3{0.9779542, -0.005750837, 0.20874014})
		renderer.UpdateCameras(renderers, &cameraMatrix)

		shadowBuffer.RenderToBuffer(func() {
			renderer.EnableDepthModeOnly()
			RenderSimulation(voxelArrayObjectRenderer)
			renderer.DisableDepthModeOnly()
		})

		// Prepare for normal rendering...
		shadowBiasMatrix := mgl32.Mat4FromCols(
			mgl32.Vec4{0.5, 0, 0, 0},
			mgl32.Vec4{0, 0.5, 0, 0},
			mgl32.Vec4{0, 0, 0.5, 0},
			mgl32.Vec4{0.5, 0.5, 0.5, 1.0})

		partialShadowMatrix := projection.Mul4(cameraMatrix.Mul4(shadowBiasMatrix))
		voxelArrayObjectRenderer.UpdateShadows(&partialShadowMatrix, shadowBuffer.GetTextureId())

		// Render the full display.
		gl.Clear(gl.COLOR_BUFFER_BIT)
		viewport.Reset()

		projection = mgl32.Perspective(mgl32.DegToRad(45.0), viewport.GetWidth()/viewport.GetHeight(), 0.1, 1000.0)
		renderer.UpdateProjections(renderers, &projection)

		cameraMatrix = camera.GetLookAtMatrix()
		renderer.UpdateCameras(renderers, &cameraMatrix)

		render()
		window.SwapBuffers()
	}
}
