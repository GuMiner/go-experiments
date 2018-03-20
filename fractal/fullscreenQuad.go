package main

// Defines a full-screen quad
import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var ccwQuad = []mgl32.Vec2{
	mgl32.Vec2{-1, 1},
	mgl32.Vec2{1, -1},
	mgl32.Vec2{1, 1},

	mgl32.Vec2{-1, -1},
	mgl32.Vec2{1, -1},
	mgl32.Vec2{-1, 1}}

type FullScreenQuad struct {
	vao uint32
	vbo uint32
}

func NewFullScreenQuad() *FullScreenQuad {
	var quad FullScreenQuad

	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &quad.vao)
	gl.BindVertexArray(quad.vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &quad.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, quad.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	// 2 -- 2 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuad)*2*4, gl.Ptr(ccwQuad), gl.STATIC_DRAW)

	return &quad
}

func (quad *FullScreenQuad) VertexCount() int32 {
	return int32(len(ccwQuad))
}

func (quad *FullScreenQuad) Render() {
	gl.BindVertexArray(quad.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, quad.VertexCount())
}

func (quad *FullScreenQuad) Delete() {
	gl.DeleteBuffers(1, &quad.vbo)
	gl.DeleteVertexArrays(1, &quad.vao)
}
