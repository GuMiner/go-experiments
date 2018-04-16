package textRenderer

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type TextProgramBuffers struct {
	vao         uint32
	positionVbo uint32
	texPosVbo   uint32
}

// Creates and returns vertex buffers to render sentence data.
func NewTextProgramBuffers() *TextProgramBuffers {
	var buffers TextProgramBuffers

	gl.GenVertexArrays(1, &buffers.vao)
	gl.BindVertexArray(buffers.vao)

	gl.EnableVertexAttribArray(0)
	gl.GenBuffers(1, &buffers.positionVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.EnableVertexAttribArray(1)
	gl.GenBuffers(1, &buffers.texPosVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers.texPosVbo)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)

	return &buffers
}

func (s *TextProgramBuffers) SendToDevice(positionBuffer []mgl32.Vec3, texturePositionBuffer []mgl32.Vec2) {
	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BindBuffer(gl.ARRAY_BUFFER, s.positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(positionBuffer)*3*4, gl.Ptr(positionBuffer), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, s.texPosVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(texturePositionBuffer)*2*4, gl.Ptr(texturePositionBuffer), gl.STATIC_DRAW)
}

func (s *TextProgramBuffers) Delete() {
	gl.DeleteBuffers(1, &s.texPosVbo)
	gl.DeleteBuffers(1, &s.positionVbo)
	gl.DeleteVertexArrays(1, &s.vao)
}
