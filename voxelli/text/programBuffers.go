package text

import "github.com/go-gl/gl/v4.5-core/gl"

type textProgramBuffers struct {
	vao         uint32
	positionVbo uint32
	texPosVbo   uint32
}

// Creates and returns vertex buffers to render sentence data.
func newTextProgramBuffers() textProgramBuffers {
	var buffers textProgramBuffers

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

	return buffers
}

func (s textProgramBuffers) Delete() {
	gl.DeleteBuffers(1, &s.texPosVbo)
	gl.DeleteBuffers(1, &s.positionVbo)
	gl.DeleteVertexArrays(1, &s.vao)
}
