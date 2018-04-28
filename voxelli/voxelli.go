package main

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.

import (
	"go-experiments/common/color"
	"go-experiments/common/config"
	"go-experiments/common/diagnostics"
	"go-experiments/common/opengl"
	"go-experiments/common/shadow"

	"go-experiments/voxelli/config"
	"go-experiments/voxelli/diagnostics"
	"go-experiments/voxelli/input"
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/text"
	"go-experiments/voxelli/vehicle"
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
	window.SetFramebufferSizeCallback(commonOpenGl.ResizeViewport)
	window.SetCursorPosCallback(input.HandleMouseMove)
	window.SetMouseButtonCallback(input.HandleMouseButton)
	window.SetKeyCallback(input.HandleKeyInput)
}

var isFpsEnabled bool = true
var isHelpTextEnabled bool = false

// Checks if the wireframe button has been toggled or not, toggling the GL setting
// This function should be called within the OpenGL update loop
func checkTextToggles() {
	if input.IsTyped(input.ToggleFpsText) {
		isFpsEnabled = !isFpsEnabled
	}

	if input.IsTyped(input.ToggleHelpText) {
		isHelpTextEnabled = !isHelpTextEnabled
	}
}

func main() {
	config.Load("./data/config.json", "./data/commonConfig.json")

	commonOpenGl.InitGlfw()
	defer glfw.Terminate()

	commonOpenGl.InitViewport()
	window, err := glfw.CreateWindow(
		int(commonOpenGl.GetWindowSize().X()),
		int(commonOpenGl.GetWindowSize().Y()),
		commonConfig.Config.Window.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	setInputCallbacks(window)
	commonOpenGl.ConfigureOpenGl()

	input.LoadKeyAssignments()
	defer input.SaveKeyAssignments()

	// Create renderers
	shadowBuffer := shadow.NewShadowBuffer()
	defer shadowBuffer.Delete()

	commonDiagnostics.InitCube()
	defer commonDiagnostics.DeleteCube()

	commonColor.InitializeColorGradient(
		commonConfig.Config.ColorGradient.Steps,
		commonConfig.Config.ColorGradient.Saturation,
		commonConfig.Config.ColorGradient.Luminosity)

	voxelArrayObjectRenderer := renderer.NewVoxelArrayObjectRenderer()
	defer voxelArrayObjectRenderer.Delete()

	textRenderer := text.NewTextRenderer(commonConfig.Config.Text.FontFile)
	defer textRenderer.Delete()

	fpsSentence := text.NewSentence(textRenderer, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	helpTextSentence := text.NewSentence(textRenderer, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 0})

	// TODO: We can't use fixed offsets because that doesn't actually work with perspective resizes
	fpsCounter := NewFpsCounter(fpsSentence, 1.0, mgl32.Vec3{-0.42, 0.33, 0.01})
	helpText := NewHelpText(helpTextSentence, mgl32.Vec3{0.18, 0.33, 0.01})

	var renderers []renderer.Renderer
	renderers = append(renderers, voxelArrayObjectRenderer)
	renderers = append(renderers, textRenderer)
	renderers = append(renderers, commonDiagnostics.GetCube())

	camera := NewCamera(
		config.Config.Camera.GetDefaultPos(),
		config.Config.Camera.GetDefaultForwards(),
		config.Config.Camera.GetDefaultUp())
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

		fpsCounter.Update(frameTime)
		opengl.CheckWireframeToggle()
		diagnostics.CheckDebugToggle()
		vehicle.CheckColorOverlayToggle()
		checkTextToggles()

		camera.Update(frameTime)

		UpdateSimulation(frameTime, elapsed)
		glfw.PollEvents()
	}

	render := func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		RenderSimulation(voxelArrayObjectRenderer)

		if isFpsEnabled {
			fpsCounter.Render(camera)
		}

		if isHelpTextEnabled {
			helpText.Render(camera)
		}
	}

	for !window.ShouldClose() {
		update()

		// Render to shadow buffer
		projection, cameraMatrix := shadowBuffer.PrepareCamera()

		renderer.UpdateProjections(renderers, &projection)
		renderer.UpdateCameras(renderers, &cameraMatrix)

		shadowBuffer.RenderToBuffer(func() {
			gl.Clear(gl.DEPTH_BUFFER_BIT)
			gl.CullFace(gl.FRONT)

			renderer.EnableDepthModeOnly()
			RenderSimulation(voxelArrayObjectRenderer)
			renderer.DisableDepthModeOnly()

			gl.CullFace(gl.BACK)
		})

		// Prepare for normal rendering...
		shadowBiasMatrix := mgl32.Mat4FromRows(
			mgl32.Vec4{0.5, 0, 0, 0.5},
			mgl32.Vec4{0, 0.5, 0, 0.5},
			mgl32.Vec4{0, 0, 0.5, 0.5},
			mgl32.Vec4{0, 0, 0, 1.0})

		partialShadowMatrix := shadowBiasMatrix.Mul4(projection.Mul4(cameraMatrix))
		voxelArrayObjectRenderer.UpdateShadows(&partialShadowMatrix, shadowBuffer.GetTextureId())

		// Render the full display.
		commonOpenGl.ResetViewport()

		projection = mgl32.Perspective(
			mgl32.DegToRad(commonConfig.Config.Perspective.FovY),
			commonOpenGl.GetWindowSize().X()/commonOpenGl.GetWindowSize().Y(),
			commonConfig.Config.Perspective.Near,
			commonConfig.Config.Perspective.Far)
		renderer.UpdateProjections(renderers, &projection)

		cameraMatrix = camera.GetLookAtMatrix()
		renderer.UpdateCameras(renderers, &cameraMatrix)

		render()
		window.SwapBuffers()
	}
}
