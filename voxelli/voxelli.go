package main

import (
	"fmt"
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

var mousePos mgl32.Vec2 = mgl32.Vec2{0, 0}

func handleMouseMove(window *glfw.Window, xPos float64, yPos float64) {
	mousePos = mgl32.Vec2{float32(xPos), float32(yPos)}
}

func setInputCallbacks(window *glfw.Window) {
	window.SetFramebufferSizeCallback(viewport.HandleResize)
	window.SetCursorPosCallback(handleMouseMove)
}

func main() {
	initGlfw()
	defer glfw.Terminate()

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

	camera := mgl32.LookAtV(mgl32.Vec3{-60, -60, 80}, mgl32.Vec3{200, 200, 00}, mgl32.Vec3{0, 0, 1})
	voxelObjectRenderer.UpdateCamera(&camera)

	startTime := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update uniforms that need updating
		elapsed := time.Since(startTime)

		// Don't distort on resize
		if !viewport.PerspectiveMatrixUpdated() {
			projection := mgl32.Perspective(mgl32.DegToRad(45.0), viewport.GetWidth()/viewport.GetHeight(), 0.1, 1000.0)
			voxelObjectRenderer.UpdateProjection(&projection)
		}

		roadwayRenderer.Render(simpleRoadway)

		// Draw a few cars
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				xCarOffset := i*(longCar.maxBounds.X()-longCar.minBounds.X()) + 4
				yCarOffset := j*(longCar.maxBounds.Y()-longCar.minBounds.Y()) + 4
				modelMatrix := mgl32.HomogRotate3D(0.5*float32(elapsed)/float32(time.Second), mgl32.Vec3{0, 0, 1})
				modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(float32(xCarOffset), float32(yCarOffset), 0.0))

				voxelObjectRenderer.Render(longCar, &modelMatrix)
			}
		}

		// for i := 0; i < 20; i++ {
		// 	 gl.Uniform1f(timeLoc, (float32(elapsed)/float32(time.Second))+float32(i)*0.3+float32(i)*0.5)
		//
		// 	 model := mgl32.Translate3D(float32(3*i), 0, float32(3*i))
		// 	 gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		//
		// 	 cube.Render()
		// }

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
