package main

// Simplifies creating a GLSL GPU shading program
import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

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
func createProgram(baseProgramName string) (program uint32) {
	vertexShader, err := compileShader(baseProgramName+".vs", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(baseProgramName+".fs", gl.FRAGMENT_SHADER)
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
