package commonOpenGl

// Simplifies creating a GLSL GPU shading program
import (
	"fmt"
	"go-experiments/common/commonio"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// Also pulled mostly from https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go

func compileShader(shaderPath string, shaderType uint32) (shader uint32) {
	source := commonIo.ReadFile(shaderPath)
	return compileShaderFromSource(source, shaderType)
}

func compileShaderFromSource(source string, shaderType uint32) (shader uint32) {
	shader = gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")

	gl.ShaderSource(shader, 1, csources, nil)
	gl.CompileShader(shader)
	free()

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("Failed to compile %v: %v", source, log))
	}

	return shader
}

// Creates an OpenGL program from a vertex shader and fragment shader, returning the program ID.
// Pulled mostly directly from https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
func CreateProgram(baseProgramName string) (program uint32) {
	vertexShader := compileShader(baseProgramName+".vs", gl.VERTEX_SHADER)
	fragmentShader := compileShader(baseProgramName+".fs", gl.FRAGMENT_SHADER)
	return createProgramFromShaders(vertexShader, fragmentShader)
}

func CreateProgramFromSource(vertexShaderSource string, fragmentShaderSource string) (program uint32) {
	vertexShader := compileShaderFromSource(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShaderFromSource(fragmentShaderSource, gl.FRAGMENT_SHADER)
	return createProgramFromShaders(vertexShader, fragmentShader)
}

func createProgramFromShaders(vertexShader uint32, fragmentShader uint32) (program uint32) {
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

		panic(fmt.Errorf("Failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}
