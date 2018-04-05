package main

import (
	"fmt"
	"go-experiments/voxelli/input"
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/text"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/viewport"
	"go-experiments/voxelli/voxel"
	"go-experiments/voxelli/voxelArray"
	"math/rand"
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

func setInputCallbacks(window *glfw.Window) {
	window.SetFramebufferSizeCallback(viewport.HandleResize)
	window.SetCursorPosCallback(input.HandleMouseMove)
	window.SetMouseButtonCallback(input.HandleMouseButton)
	window.SetKeyCallback(input.HandleKeyInput)
}

var wasDebugPressed bool = false
var isDebug bool = false

func checkDebugToggle() {
	if !wasDebugPressed && input.PressedKeys[glfw.KeyT] {
		wasDebugPressed = true
		isDebug = !isDebug
	}

	if wasDebugPressed && !input.PressedKeys[glfw.KeyT] {
		wasDebugPressed = false
	}
}

func debugDrawCarInfo(debugCube *Cube, car *vehicle.Vehicle, elapsed float32, boundaries []float32) {
	eyePositions, eyeDirections := car.GetEyes()

	// Debug draw where we are looking, assuming two eyes only
	model := mgl32.Translate3D(eyePositions[0].X()+2*eyeDirections[0].X(), eyePositions[0].Y()+2*eyeDirections[0].Y(), 8)
	debugCube.Render(elapsed, mgl32.Vec4{0.0, 1.0, 0.0, 1.0}, &model)

	model = mgl32.Translate3D(eyePositions[1].X()+2*eyeDirections[1].X(), eyePositions[1].Y()+2*eyeDirections[1].Y(), 8)
	debugCube.Render(elapsed, mgl32.Vec4{0.0, 1.0, 1.0, 1.0}, &model) // Cyan

	// Debug draw where the boundaries are.
	eyePositions[0] = eyePositions[0].Add(eyeDirections[0].Mul(boundaries[0]))
	eyePositions[1] = eyePositions[1].Add(eyeDirections[1].Mul(boundaries[1]))

	model = mgl32.Translate3D(eyePositions[0].X(), eyePositions[0].Y(), 8)
	debugCube.Render(elapsed, mgl32.Vec4{1.0, 0.0, 0.0, 1.0}, &model)

	model = mgl32.Translate3D(eyePositions[1].X(), eyePositions[1].Y(), 8)
	debugCube.Render(elapsed, mgl32.Vec4{1.0, 1.0, 0.0, 1.0}, &model) // Yellow
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

	// Create renderers
	debugCube := NewCube()
	defer debugCube.Delete()

	voxelArrayObjectRenderer := renderer.NewVoxelArrayObjectRenderer()
	defer voxelArrayObjectRenderer.Delete()

	textRenderer := text.NewTextRenderer("./data/font/DejaVuSans.ttf")
	defer textRenderer.Delete()

	var renderers []renderer.Renderer
	renderers = append(renderers, voxelArrayObjectRenderer)
	renderers = append(renderers, textRenderer)
	renderers = append(renderers, debugCube)

	// Create roadway
	simpleRoadway := NewRoadway("./data/roadways/straight_with_s-curve.txt")
	fmt.Printf("Straight roadway size: [%v, %v]\n\n", len(simpleRoadway.roadElements), len(simpleRoadway.roadElements[0]))

	roadwayDisplayer := NewRoadwayDisplayer(voxelArrayObjectRenderer)
	defer roadwayDisplayer.Delete()

	carRaw := voxel.NewVoxelObject("./data/models/car.vox")
	fmt.Printf("Vehicle objects: %v\n", len(carRaw.SubObjects))

	carModel := voxelArray.NewVoxelArrayObject(carRaw)
	defer carModel.Delete()
	fmt.Printf("Optimized Vehicle vertices: %v\n\n", carModel.Vertices)

	var cars []*vehicle.Vehicle
	for i := 0; i < 10; i++ {
		car := vehicle.NewVehicle(i, carModel)
		car.Position[0] = 10
		car.Position[1] = 10

		car.AccelPos = rand.Float32()*2 - 1
		car.SteeringPos = rand.Float32()*2 - 1

		cars = append(cars, car)
	}

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

		if input.AnyEvent() {
			opengl.CheckWireframeToggle()
			checkDebugToggle()
		}

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

		var maxScore float32 = 0.0
		maxScoreIdx := 0
		for i, car := range cars {
			oldPosition, oldOrientation := car.Update(frameTime)
			if !simpleRoadway.InAllBounds(car.GetBounds()) {
				car.AccelPos = rand.Float32()*2 - 1
				car.SteeringPos = rand.Float32()*2 - 1
				fmt.Printf("Acc: %v. Steer: %v Pos: %v\n", car.AccelPos, car.SteeringPos, car.Position)

				car.Velocity = 0
				car.Position = oldPosition
				car.Orientation = oldOrientation
				car.Score = 0
			}

			if car.Score > maxScore {
				maxScoreIdx = i
				maxScore = car.Score
			}

			eyePositions, eyeDirections := car.GetEyes()
			boundaryLengths := simpleRoadway.GetBoundaries(eyePositions, eyeDirections)
			if isDebug {
				debugDrawCarInfo(debugCube, car, elapsed, boundaryLengths)
			}
		}

		for i, car := range cars {
			emphasize := i == maxScoreIdx
			car.Render(emphasize, voxelArrayObjectRenderer)
		}

		ident := mgl32.Ident4()
		textRenderer.Render("Hello World!", &ident)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
