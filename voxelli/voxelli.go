package main

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.

import (
	"fmt"
	"go-experiments/voxelli/diagnostics"
	"go-experiments/voxelli/genetics"
	"go-experiments/voxelli/input"
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/text"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/viewport"
	"go-experiments/voxelli/voxel"
	"go-experiments/voxelli/voxelArray"
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

func debugDrawCarInfo(car *vehicle.Vehicle, elapsed float32, boundaries []float32) {
	eyePositions, eyeDirections := car.GetEyes()

	// Debug draw where we are looking, assuming two eyes only
	model := mgl32.Translate3D(eyePositions[0].X()+2*eyeDirections[0].X(), eyePositions[0].Y()+2*eyeDirections[0].Y(), 8)
	diagnostics.Render(mgl32.Vec4{0.0, 1.0, 0.0, 1.0}, &model)

	model = mgl32.Translate3D(eyePositions[1].X()+2*eyeDirections[1].X(), eyePositions[1].Y()+2*eyeDirections[1].Y(), 8)
	diagnostics.Render(mgl32.Vec4{0.0, 1.0, 1.0, 1.0}, &model) // Cyan

	model = mgl32.Translate3D(eyePositions[2].X()+2*eyeDirections[2].X(), eyePositions[2].Y()+2*eyeDirections[2].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 1.0, 1.0, 1.0}, &model) // Cyan

	// Debug draw where the boundaries are.
	eyePositions[0] = eyePositions[0].Add(eyeDirections[0].Mul(boundaries[0]))
	eyePositions[1] = eyePositions[1].Add(eyeDirections[1].Mul(boundaries[1]))
	eyePositions[2] = eyePositions[2].Add(eyeDirections[2].Mul(boundaries[2]))

	model = mgl32.Translate3D(eyePositions[0].X(), eyePositions[0].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 0.0, 0.0, 1.0}, &model)

	model = mgl32.Translate3D(eyePositions[1].X(), eyePositions[1].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 1.0, 0.0, 1.0}, &model) // Yellow

	model = mgl32.Translate3D(eyePositions[2].X(), eyePositions[2].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 1.0, 1.0, 1.0}, &model) // Yellow
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
	diagnostics.InitCube()
	defer diagnostics.DeleteCube()

	voxelArrayObjectRenderer := renderer.NewVoxelArrayObjectRenderer()
	defer voxelArrayObjectRenderer.Delete()

	textRenderer := text.NewTextRenderer("./data/font/DejaVuSans.ttf")
	defer textRenderer.Delete()

	var renderers []renderer.Renderer
	renderers = append(renderers, voxelArrayObjectRenderer)
	renderers = append(renderers, textRenderer)
	renderers = append(renderers, diagnostics.GetCube())

	// Create roadway
	simpleRoadway := roadway.NewRoadway("./data/roadways/straight_with_s-curve.txt")

	roadwayDisplayer := roadway.NewRoadwayDisplayer(voxelArrayObjectRenderer)
	defer roadwayDisplayer.Delete()

	// Create car
	carVoxelObject := voxel.NewVoxelObject("./data/models/car.vox")
	fmt.Printf("Vehicle objects: %v\n", len(carVoxelObject.SubObjects))

	carModel := voxelArray.NewVoxelArrayObject(carVoxelObject)
	defer carModel.Delete()
	fmt.Printf("Optimized Vehicle vertices: %v\n\n", carModel.Vertices)

	// Create agent population
	agentEvolver := genetics.NewPopulation(100, func(id int) *genetics.Agent {
		return genetics.NewAgent(id, carModel, mgl32.Vec2{10, 10})
	})

	camera := NewCamera(mgl32.Vec3{140, 300, 300}, mgl32.Vec3{-1, 0, 0}, mgl32.Vec3{0, 0, 1})
	defer camera.CachePosition()

	cameraMatrix := camera.GetLookAtMatrix()
	renderer.UpdateCameras(renderers, &cameraMatrix)

	startTime := time.Now()
	lastElapsed := float32(0.0)
	elapsed := lastElapsed
	for !window.ShouldClose() {
		lastElapsed = elapsed
		elapsed = float32(time.Since(startTime)) / float32(time.Second)
		frameTime := elapsed - lastElapsed

		// Start rendering and updating
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		opengl.CheckWireframeToggle()
		diagnostics.CheckDebugToggle()

		// Update our camera if we have motion
		if camera.Update(frameTime, &cameraMatrix) {
			renderer.UpdateCameras(renderers, &cameraMatrix)
		}

		// Don't distort on resize
		if !viewport.PerspectiveMatrixUpdated() {
			projection := mgl32.Perspective(mgl32.DegToRad(45.0), viewport.GetWidth()/viewport.GetHeight(), 0.1, 1000.0)
			renderer.UpdateProjections(renderers, &projection)
		}

		roadwayDisplayer.Render(simpleRoadway)

		// Update and render the agents
		agentEvolver.Update(frameTime, func(agent *genetics.Agent) {
			agent.Update(frameTime, simpleRoadway)

			if diagnostics.IsDebug() {
				eyePositions, eyeDirections := agent.GetCar().GetEyes()
				boundaryLengths := simpleRoadway.GetBoundaries(eyePositions, eyeDirections)
				debugDrawCarInfo(agent.GetCar(), elapsed, boundaryLengths)
			}
		})

		agentEvolver.Render(func(agent *genetics.Agent) {
			agent.Render(voxelArrayObjectRenderer)
		})

		moveOver := mgl32.Translate3D(200, 200, 50)
		textRenderer.Render("Hello World!", &moveOver)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
