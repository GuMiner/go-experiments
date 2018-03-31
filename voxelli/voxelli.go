package main

import (
	"fmt"
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
	voxelArrayObjectRenderer := renderer.NewVoxelArrayObjectRenderer()
	defer voxelArrayObjectRenderer.Delete()

	textRenderer := text.NewTextRenderer("./data/font/DejaVuSans.ttf")
	defer textRenderer.Delete()

	var renderers []renderer.Renderer
	renderers = append(renderers, voxelArrayObjectRenderer)
	renderers = append(renderers, textRenderer)

	// Create roadway
	simpleRoadway := NewRoadway("./data/roadways/straight_with_s-curve.txt")
	fmt.Printf("Straight roadway size: [%v, %v]\n\n", len(simpleRoadway.roadElements), len(simpleRoadway.roadElements[0]))

	roadwayDisplayer := NewRoadwayDisplayer(voxelArrayObjectRenderer)
	defer roadwayDisplayer.Delete()

	car := vehicle.NewVehicle("./data/models/car.vox")
	defer car.Delete()

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

		// Draw a few cars
		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				xCarOffset := float32(i)*car.HalfSize.X()*2 + 4
				yCarOffset := float32(j)*car.HalfSize.Y()*2 + 4
				rotateMatrix := mgl32.HomogRotate3D(0.5*elapsed, mgl32.Vec3{0, 0, 1})

				translateMatrix := mgl32.Translate3D(xCarOffset, yCarOffset, 1)
				modelMatrix := rotateMatrix.Mul4(translateMatrix)

				// TODO Remove very soon
				voxelArrayObjectRenderer.Render(car.Shape, &modelMatrix)
			}
		}

		ident := mgl32.Ident4()
		textRenderer.Render("Hello World!", &ident)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
