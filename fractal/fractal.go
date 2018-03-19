package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"time"

	color "github.com/gerow/go-color"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func setup() (vao uint32, vbo uint32, vertexCount int32) {
	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	ccwQuad := []mgl32.Vec2{
		mgl32.Vec2{-1, 1},
		mgl32.Vec2{1, -1},
		mgl32.Vec2{1, 1},

		mgl32.Vec2{-1, -1},
		mgl32.Vec2{1, -1},
		mgl32.Vec2{-1, 1}}

	// 2 -- 2 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuad)*2*4, gl.Ptr(ccwQuad), gl.STATIC_DRAW)

	return vao, vbo, int32(len(ccwQuad))
}

func teardown(vao uint32, vbo uint32) {
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteVertexArrays(1, &vao)
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

	// TODO: Re-enable post-debug to get a performance boost
	// gl.Disable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

func readFile(path string) string {
	fileAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(fileAsBytes)
}

// Also pulled mostly from https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
func compileShader(shaderPath string, shaderType uint32) (shader uint32, errorIfAny error) {
	shader = gl.CreateShader(shaderType)

	source := readFile(shaderPath)

	csources, free := gl.Strs(source)

	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// Pulled mostly directly from https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
func createProgram() (program uint32) {
	vertexShader, err := compileShader("./shaders/juliaFractal.vs", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader("./shaders/juliaFractal.fs", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program = gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func createColorGradient(maxIterations int32, saturation float32, luminosity float32) []mgl32.Vec3 {
	colorGradient := make([]mgl32.Vec3, maxIterations)

	for idx := range colorGradient {
		hue := float32(idx) / float32(maxIterations)
		color := color.HSL{float64(hue), float64(saturation), float64(luminosity)}.ToRGB()
		colorGradient[idx] = mgl32.Vec3{float32(color.R), float32(color.G), float32(color.B)}
	}

	return colorGradient
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	configureOpenGl()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version: ", version)

	program := createProgram()

	mouseUniformLoc := gl.GetUniformLocation(program, gl.Str("c\x00"))
	timeLoc := gl.GetUniformLocation(program, gl.Str("time\x00"))
	fractalGradientLoc := gl.GetUniformLocation(program, gl.Str("fractalGradient\x00"))
	maxIterationsLoc := gl.GetUniformLocation(program, gl.Str("maxIterations\x00"))

	gl.ActiveTexture(gl.TEXTURE0)
	var fractalGradientTextureID uint32

	var maxIterations int32 = 1000
	colorGradient := createColorGradient(maxIterations, 1.0, 0.5)

	gl.GenTextures(1, &fractalGradientTextureID)
	gl.BindTexture(gl.TEXTURE_1D, fractalGradientTextureID)
	gl.TexStorage1D(gl.TEXTURE_1D, 1, gl.RGB32F, maxIterations)
	gl.TexSubImage1D(gl.TEXTURE_1D, 0, 0, maxIterations, gl.RGB, gl.FLOAT, gl.Ptr(colorGradient))

	vao, vbo, vertexCount := setup()

	startTime := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_1D, fractalGradientTextureID)
		gl.Uniform1i(fractalGradientLoc, gl.TEXTURE0-gl.TEXTURE0) // Yes it is zero.

		gl.Uniform1i(maxIterationsLoc, maxIterations)

		elapsed := time.Since(startTime)
		gl.Uniform1f(timeLoc, float32(elapsed)/float32(time.Second))

		// TODO: Replace with the mouse positionings
		gl.Uniform2f(mouseUniformLoc, 0.0, 0.0)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteTextures(1, &fractalGradientTextureID)
	teardown(vao, vbo)
	gl.DeleteProgram(program)
}
