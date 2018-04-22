package main

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.

import (
	"go-experiments/common/opengl"

	"go-experiments/voxelli/color"
	"go-experiments/voxelli/config"
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
	config.Load("./data/config.json")

	commonOpenGl.InitGlfw()
	defer glfw.Terminate()

	viewport.Init()
	window, err := glfw.CreateWindow(
		int(viewport.GetWidth()),
		int(viewport.GetHeight()),
		config.Config.Window.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	setInputCallbacks(window)
	commonOpenGl.ConfigureOpenGl()

	input.LoadKeyAssignments()
	defer input.SaveKeyAssignments()

	// Create renderers
	shadowBuffer := NewShadowBuffer()
	defer shadowBuffer.Delete()

	diagnostics.InitCube()
	defer diagnostics.DeleteCube()

	color.InitializeColorGradient(
		config.Config.ColorGradient.Steps,
		config.Config.ColorGradient.Saturation,
		config.Config.ColorGradient.Luminosity)

	voxelArrayObjectRenderer := renderer.NewVoxelArrayObjectRenderer()
	defer voxelArrayObjectRenderer.Delete()

	textRenderer := text.NewTextRenderer(config.Config.Text.FontFile)
	defer textRenderer.Delete()

	fpsSentence := text.NewSentence(textRenderer, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	helpTextSentence := text.NewSentence(textRenderer, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 0})

	// TODO: We can't use fixed offsets because that doesn't actually work with perspective resizes
	fpsCounter := NewFpsCounter(fpsSentence, 1.0, mgl32.Vec3{-0.42, 0.33, 0.01})
	helpText := NewHelpText(helpTextSentence, mgl32.Vec3{0.18, 0.33, 0.01})

	var renderers []renderer.Renderer
	renderers = append(renderers, voxelArrayObjectRenderer)
	renderers = append(renderers, textRenderer)
	renderers = append(renderers, diagnostics.GetCube())

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
		gl.Viewport(0, 0, shadowBuffer.Width, shadowBuffer.Height)

		// Hardcoded for the following properties:
		// 1. Good shadow angle
		// 2. Good usage of the depth buffer range
		projection := mgl32.Ortho(
			config.Config.Shadows.Projection.Left,
			config.Config.Shadows.Projection.Right,
			config.Config.Shadows.Projection.Bottom,
			config.Config.Shadows.Projection.Top,
			config.Config.Shadows.Projection.Near,
			config.Config.Shadows.Projection.Far)
		renderer.UpdateProjections(renderers, &projection)

		position := mgl32.Vec3{
			config.Config.Shadows.Position.X,
			config.Config.Shadows.Position.Y,
			config.Config.Shadows.Position.Z}
		cameraMatrix := mgl32.LookAtV(
			position,
			position.Add(mgl32.Vec3{
				config.Config.Shadows.Forwards.X,
				config.Config.Shadows.Forwards.Y,
				config.Config.Shadows.Forwards.Z}),
			mgl32.Vec3{
				config.Config.Shadows.Up.X,
				config.Config.Shadows.Up.Y,
				config.Config.Shadows.Up.Z})
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
		viewport.Reset()

		projection = mgl32.Perspective(
			mgl32.DegToRad(config.Config.Perspective.FovY),
			viewport.GetWidth()/viewport.GetHeight(),
			config.Config.Perspective.Near,
			config.Config.Perspective.Far)
		renderer.UpdateProjections(renderers, &projection)

		cameraMatrix = camera.GetLookAtMatrix()
		renderer.UpdateCameras(renderers, &cameraMatrix)

		render()
		window.SwapBuffers()
	}
}
