package text

import (
	"go-experiments/voxelli/opengl"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type textRendererProgram struct {
	shaderProgram uint32
	projectionLoc int32
	cameraLoc     int32
	modelLoc      int32
	fontImageLoc  int32

	foregroundColorLoc int32
	backgroundColorLoc int32
}

func newTextRendererProgram() textRendererProgram {
	var program textRendererProgram
	program.shaderProgram = opengl.CreateProgram("./text/textRenderer")

	program.projectionLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("projection\x00"))
	program.cameraLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("camera\x00"))
	program.modelLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("model\x00"))
	program.fontImageLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("fontImage\x00"))

	program.foregroundColorLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("foregroundColor\x00"))
	program.backgroundColorLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("backgroundColor\x00"))

	return program
}

func (r *textRendererProgram) Delete() {
	gl.DeleteProgram(r.shaderProgram)
}
