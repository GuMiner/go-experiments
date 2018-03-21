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

	longCar := NewVoxelObject("./models/long_car.vox")
	fmt.Printf("Long Car objects: %v\n", len(longCar.subObjects))

	straightRoad := NewVoxelObject("./models/road_straight.vox")
	fmt.Printf("Straight Road objects: %v\n", len(straightRoad.subObjects))

	program := createProgram("./shaders/basicRenderer")
	defer gl.DeleteProgram(program)

	// Get locations of everything used in this program.
	projectionLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraLoc := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))
	timeLoc := gl.GetUniformLocation(program, gl.Str("runTime\x00"))

	cube := NewCube()
	defer cube.Delete()

	// Enable our program and do first-time uniform setup
	gl.UseProgram(program)

	camera := mgl32.LookAtV(mgl32.Vec3{-3, 4, -3}, mgl32.Vec3{20, 0, 20}, mgl32.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(cameraLoc, 1, false, &camera[0])

	startTime := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update uniforms that need updating
		elapsed := time.Since(startTime)

		// Don't distort on resize
		if !viewport.PerspectiveMatrixUpdated() {
			projection := mgl32.Perspective(mgl32.DegToRad(45.0), viewport.GetWidth()/viewport.GetHeight(), 0.1, 1000.0)
			gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])
		}

		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {

				gl.Uniform1f(timeLoc, (float32(elapsed)/float32(time.Second))+float32(i)*0.3+float32(j)*0.5)

				model := mgl32.Translate3D(float32(3*i), 0, float32(3*j))
				gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])

				cube.Render()
			}
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
