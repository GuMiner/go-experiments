package main

// Defines a small cube
import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Directly from https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
var cubeVertices = []mgl32.Vec3{
	// Bottom
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},

	// Top
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},

	// Front
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},

	// Back
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},

	// Left
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},

	// Right
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, 0.5}}

var cubeColorVertices = []mgl32.Vec3{
	// Bottom
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.5, 0.2},

	// Top
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.2},

	// Front
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.2},

	// Back
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.5},

	// Left
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.5},

	// Right
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.2}}

type Cube struct {
	vao         uint32
	positionVbo uint32
	colorVbo    uint32
}

func NewCube() *Cube {
	var cube Cube

	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &cube.vao)
	gl.BindVertexArray(cube.vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &cube.positionVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*3*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(1)
	gl.GenBuffers(1, &cube.colorVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.colorVbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.BufferData(gl.ARRAY_BUFFER, len(cubeColorVertices)*3*4, gl.Ptr(cubeColorVertices), gl.STATIC_DRAW)

	return &cube
}

func (cube *Cube) VertexCount() int32 {
	return int32(len(cubeVertices))
}

func (cube *Cube) Render() {
	gl.BindVertexArray(cube.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, cube.VertexCount())
}

func (cube *Cube) Delete() {
	gl.DeleteBuffers(1, &cube.positionVbo)
	gl.DeleteBuffers(1, &cube.colorVbo)
	gl.DeleteVertexArrays(1, &cube.vao)
}
