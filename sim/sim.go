package main

import (
	"go-experiments/common/commoncolor"
	"go-experiments/common/commonconfig"
	"go-experiments/common/commonopengl"

	"github.com/go-gl/mathgl/mgl32"

	"go-experiments/sim/config"
	"go-experiments/sim/engine"
	"go-experiments/sim/input"
	"go-experiments/sim/input/editorEngine"
	"go-experiments/sim/ui"
	"go-experiments/sim/ui/flat"

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
	config.Load("./data/config/", "./data/commonConfig.json")

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

	paused := false
	update := func() {
		lastElapsed = elapsed
		elapsed = float32(time.Since(startTime)) / float32(time.Second)
		frameTime = elapsed - lastElapsed

		// Must be first.
		glfw.PollEvents()
		mouseMoved := input.MouseMoveEvent
		input.MouseMoveEvent = false

		camera.Update(frameTime)

		editorStateUpdated, editorSubStateUpdated := editorEngine.Update()
		if editorStateUpdated || editorSubStateUpdated {
			// The edit state has updated, update as needed
			ui.UpdateEditorState(editorEngine.EngineState, window)
		}

		// Load new terrain regions based on what is visible.
		engine.PrecacheRegions(camera.ComputePrecacheRegions())

		visibleRegions := camera.ComputeVisibleRegions()
		for _, region := range visibleRegions {
			subMap := engine.GetRegionMap(region)

			overlay, isNew := terrainOverlays.GetOrAddTerrainOverlay(region.X(), region.Y())
			if isNew || subMap.Dirty {
				overlay.SetTerrain(subMap.Texels)
				subMap.Dirty = false
			}
		}

		boardPos := camera.MapPixelPosToBoard(input.MousePos)
		if editorStateUpdated || mouseMoved {
			engine.Hypotheticals.ComputeHypotheticalRegion(engine, &editorEngine.EngineState)
			engine.ComputeSnapNodes(&editorEngine.EngineState)
		}

		if input.MousePressEvent {
			engine.MousePress(boardPos, editorEngine.EngineState)
			input.MousePressEvent = false
		}

		if input.MouseReleaseEvent {
			engine.MouseRelease(boardPos, editorEngine.EngineState)
			input.MouseReleaseEvent = false
		}

		if mouseMoved {
			engine.MouseMoved(boardPos, editorEngine.EngineState)
		}

		if input.IsTyped(input.CancelKey) {
			engine.CancelState(editorEngine.EngineState)
		}

		if input.IsTyped(input.PauseKey) {
			paused = !paused
		}

		if !paused {
			engine.StepSim(frameTime)
		}
		engine.StepEdit(frameTime, editorEngine.EngineState)
	}

	render := func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render each visible region
		visibleRegions := camera.ComputeVisibleRegions()
		ui.Ui.OverlayProgram.PreRender()
		for _, region := range visibleRegions {
			overlay, _ := terrainOverlays.GetOrAddTerrainOverlay(region.X(), region.Y())
			overlay.UpdateCameraOffset(region.X(), region.Y(), camera)
			ui.Ui.OverlayProgram.Render(overlay.GetOverlay())
		}

		ui.Ui.RegionProgram.PreRender()
		for _, hypotheticalRegion := range engine.Hypotheticals.Regions {
			mappedRegion := camera.MapEngineRegionToScreen(&hypotheticalRegion.Region)
			ui.Ui.RegionProgram.Render(mappedRegion, hypotheticalRegion.Color)
		}

		ui.Ui.LinesProgram.PreRender()
		for _, hypotheticalLine := range engine.Hypotheticals.Lines {
			mappedLine := camera.MapEngineLineToScreen(hypotheticalLine.Line)
			ui.Ui.LinesProgram.Render([][2]mgl32.Vec2{mappedLine}, hypotheticalLine.Color)
		}

		flat.RenderPowerPlants(engine.GetPowerGrid(), camera, ui.Ui.RegionProgram)
		flat.RenderPowerLines(engine.GetPowerGrid(), camera, ui.Ui.LinesProgram)

		flat.RenderRoadLines(engine.GetRoadGrid(), camera, ui.Ui.LinesProgram)

		flat.RenderSnapNodes(engine.GetSnapElements(), camera, ui.Ui.RegionProgram)
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
