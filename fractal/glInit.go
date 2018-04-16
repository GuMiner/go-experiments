package main

// Simplifies OpenGL initialization
import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func logOpenGlInfo() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version: ", version)
}

func configureOpenGl() {
	// Startup OpenGL bindings
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	gl.ClearColor(0.2, 0.5, 1.0, 1.0)

	glfw.SwapInterval(1)

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.LINE_SMOOTH)
	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	logOpenGlInfo()
}
