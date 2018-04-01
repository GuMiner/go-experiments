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

	car := vehicle.NewVehicle(carModel)
	car.Position[0] = 10
	car.Position[1] = 10

	car.AccelPos = rand.Float32()*2 - 1
	car.SteeringPos = rand.Float32()*2 - 1

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

		oldPosition, oldOrientation := car.Update(frameTime)
		if !simpleRoadway.InAllBounds(car.GetBounds()) {
			car.AccelPos = rand.Float32()*2 - 1
			car.SteeringPos = rand.Float32()*2 - 1
			fmt.Printf("Acc: %v. Steer: %v Pos: %v\n", car.AccelPos, car.SteeringPos, car.Position)

			car.Velocity = 0
			car.Position = oldPosition
			car.Orientation = oldOrientation
		}

		// var height float32 = 6.0
		// bounds := car.GetBounds()
		// for _, bound := range bounds {
		// 	color := mgl32.Vec4{0.0, 1.0, 0.0, 1.0}
		// 	model := mgl32.Translate3D(bound.X(), bound.Y(), height)
		// 	debugCube.Render(0, color, &model)
		// }
		//
		// pcolor := mgl32.Vec4{0.0, 1.0, 1.0, 1.0}
		// pmodel := mgl32.Translate3D(car.Position.X(), car.Position.Y(), height)
		// debugCube.Render(0, pcolor, &pmodel)

		// debugStep := float32(GetGridSize()*6) / 80
		// for i := 0; i < 80; i++ {
		// 	for j := 0; j < 80; j++ {
		// 		height := 5
		// 		color := mgl32.Vec4{1.0, 1.0, 0.0, 1.0}
		// 		if simpleRoadway.InBounds(mgl32.Vec2{float32(i) * debugStep, float32(j) * debugStep}) {
		// 			height += 5
		// 			color[0] = 0
		// 		}
		//
		// 		model := mgl32.Translate3D(float32(i)*debugStep, float32(j)*debugStep, float32(height))
		// 		debugCube.Render(0, color, &model)
		// 	}
		// }

		car.Render(voxelArrayObjectRenderer)

		ident := mgl32.Ident4()
		textRenderer.Render("Hello World!", &ident)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
