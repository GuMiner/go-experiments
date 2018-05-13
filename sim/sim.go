package main

import (
	"go-experiments/common/color"
	"go-experiments/common/config"
	"go-experiments/common/opengl"

	"go-experiments/sim/config"
	"go-experiments/sim/engine"
	"go-experiments/sim/input"
	"go-experiments/sim/input/editorEngine"
	"go-experiments/sim/visuals/ui"
	"go-experiments/sim/visuals/ui/flat"

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

	input.CreateDefaultKeyMap()
	window.SetCursorPosCallback(input.HandleMouseMove)
	window.SetMouseButtonCallback(input.HandleMouseButton)
	window.SetScrollCallback(input.HandleMouseScroll)
	window.SetKeyCallback(input.HandleKeyInput)
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

	commonColor.InitializeColorGradient(
		commonConfig.Config.ColorGradient.Steps,
		commonConfig.Config.ColorGradient.Saturation,
		commonConfig.Config.ColorGradient.Luminosity)

	ui.Init(window)
	defer ui.Delete()

	overlayProgram := ui.NewOverlayShaderProgram()
	defer overlayProgram.Delete()

	editorEngine.Init()

	camera := flat.NewCamera()

	// Setup simulation
	engine := engine.NewEngine()

	terrainOverlays := flat.NewTerrainOverlayManager()
	defer terrainOverlays.Delete()

	startTime := time.Now()
	frameTime := float32(0.1)
	lastElapsed := float32(0.0)
	elapsed := lastElapsed

	update := func() {
		lastElapsed = elapsed
		elapsed = float32(time.Since(startTime)) / float32(time.Second)
		frameTime = elapsed - lastElapsed

		// Must be first.
		glfw.PollEvents()

		camera.Update(frameTime)

		editorStateUpdated, _ := editorEngine.Update()
		if editorStateUpdated {
			// The edit state has updated, update as needed
			ui.UpdateEditorState(editorEngine.EngineState, window)
		}

		// Load new terrain regions based on what is visible.
		engine.PrecacheRegions(camera.ComputePrecacheRegions())

		visibleRegions := camera.ComputeVisibleRegions()
		for _, region := range visibleRegions {
			subMap := engine.GetRegionMap(region)

			overlay, isNew := terrainOverlays.GetOrAddTerrainOverlay(region.X(), region.Y())
			if isNew {
				overlay.SetTerrain(subMap.Texels)
			}
		}

		boardPos := camera.MapToBoard(input.MousePos)
		if engine.HasHypotheticalRegion(boardPos, editorEngine.EngineState) {
			// isValid, region := engine.GetHypotheticalRegion(boardPos, editorEngine.EngineState)
			// TODO: Use hypothetical region for drawing what is valid or not.
		}

		if input.MousePressEvent {
			engine.MousePress(boardPos, editorEngine.EngineState)
			input.MousePressEvent = false
		}

		if input.MouseReleaseEvent {
			engine.MouseRelease(boardPos, editorEngine.EngineState)
			input.MouseReleaseEvent = false
		}
	}

	render := func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render each visible region
		visibleRegions := camera.ComputeVisibleRegions()
		overlayProgram.PreRender()
		for _, region := range visibleRegions {
			overlay, _ := terrainOverlays.GetOrAddTerrainOverlay(region.X(), region.Y())
			overlay.UpdateCameraOffset(region.X(), region.Y(), camera)
			overlayProgram.Render(overlay.GetOverlay())
		}
	}

	RenderLoop(update, render, window)
}

func RenderLoop(update, render func(), window *glfw.Window) {
	for !window.ShouldClose() {
		update()

		// Render the full display.
		commonOpenGl.ResetViewport()

		render()
		window.SwapBuffers()
	}
}
