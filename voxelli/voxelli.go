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
	colorOverrideLoc := gl.GetUniformLocation(program, gl.Str("colorOverride\x00"))

	cube := NewCube()
	defer cube.Delete()

	// Enable our program and do first-time uniform setup
	gl.UseProgram(program)

	camera := mgl32.LookAtV(mgl32.Vec3{-60, -60, 80}, mgl32.Vec3{200, 200, 00}, mgl32.Vec3{0, 0, 1})
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

		// Draw a road piece
		for _, subObject := range straightRoad.subObjects {
			for _, voxel := range subObject.voxels {
				colorVector := longCar.palette.colors[voxel.colorIdx-1].AsFloatVector()
				gl.Uniform4fv(colorOverrideLoc, 1, &colorVector[0])

				model := mgl32.Translate3D(float32(voxel.position.X()*2), float32(voxel.position.Y()*2), float32(voxel.position.Z()*2))
				gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])

				cube.Render()
			}
		}

		// Draw a few cars
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				for _, subObject := range longCar.subObjects {
					for _, voxel := range subObject.voxels {
						gl.Uniform1f(timeLoc, (float32(elapsed)/float32(time.Second))+float32(voxel.colorIdx)*0.3)

						colorVector := longCar.palette.colors[voxel.colorIdx-1].AsFloatVector()
						gl.Uniform4fv(colorOverrideLoc, 1, &colorVector[0])

						xCarOffset := i*2*(longCar.maxBounds.X()-longCar.minBounds.X()) + 4
						yCarOffset := j*2*(longCar.maxBounds.Y()-longCar.minBounds.Y()) + 4

						model := mgl32.HomogRotate3D(0.5*float32(elapsed)/float32(time.Second), mgl32.Vec3{0, 0, 1})
						model = model.Mul4(mgl32.Translate3D(float32(voxel.position.X()*2+xCarOffset), float32(voxel.position.Y()*2+yCarOffset), float32(voxel.position.Z()*2)))

						gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])

						cube.Render()
					}
				}
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
