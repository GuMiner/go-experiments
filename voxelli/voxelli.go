package main

import (
	"fmt"
	"go-experiments/voxelli/input"
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

func initGlfw() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
}

func setInputCallbacks(window *glfw.Window) {
	window.SetFramebufferSizeCallback(viewport.HandleResize)
	window.SetCursorPosCallback(input.HandleMouseMove)
	window.SetKeyCallback(input.HandleKeyInput)
}

func main() {
	initGlfw()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(int(viewport.GetWidth()), int(viewport.GetHeight()), "Voxelli", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	setInputCallbacks(window)
	configureOpenGl()

	longCar := NewVoxelObject("./data/models/long_car.vox")
	fmt.Printf("Long Car objects: %v\n", len(longCar.subObjects))

	simpleRoadway := NewRoadway("./data/roadways/straight_with_s-curve.txt")
	fmt.Printf("Straight roadway size: [%v, %v]\n", len(simpleRoadway.roadElements), len(simpleRoadway.roadElements[0]))

	voxelObjectRenderer := NewVoxelObjectRenderer()
	defer voxelObjectRenderer.Delete()

	roadwayRenderer := NewRoadwayRenderer(voxelObjectRenderer)

	camera := NewCamera(mgl32.Vec3{140, 300, 300}, mgl32.Vec3{-1, 0, -1}, mgl32.Vec3{0, 0, 1})
	cameraMatrix := camera.GetLookAtMatrix()
	voxelObjectRenderer.UpdateCamera(&cameraMatrix)

	startTime := time.Now()
	lastElapsed := float32(0.0)
	elapsed := lastElapsed
	for !window.ShouldClose() {
		lastElapsed = elapsed
		elapsed = float32(time.Since(startTime)) / float32(time.Second)
		frameTime := elapsed - lastElapsed

		// Start rendering and updating
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update our camera if we have motion
		if camera.Update(frameTime, &cameraMatrix) {
			voxelObjectRenderer.UpdateCamera(&cameraMatrix)
		}

		// Don't distort on resize
		if !viewport.PerspectiveMatrixUpdated() {
			projection := mgl32.Perspective(mgl32.DegToRad(45.0), viewport.GetWidth()/viewport.GetHeight(), 0.1, 1000.0)
			voxelObjectRenderer.UpdateProjection(&projection)
		}

		roadwayRenderer.Render(simpleRoadway)

		// Draw a few cars
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				xCarOffset := i*(longCar.maxBounds.X()-longCar.minBounds.X()) + 4
				zCarOffset := j*(longCar.maxBounds.Z()-longCar.minBounds.Z()) + 4
				rotateMatrix := mgl32.HomogRotate3D(0.5*elapsed, mgl32.Vec3{0, 0, 1})
				translateMatrix := mgl32.Translate3D(float32(xCarOffset), 0.0, float32(zCarOffset))
				modelMatrix := rotateMatrix.Mul4(translateMatrix)

				voxelObjectRenderer.Render(longCar, &modelMatrix)
			}
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
