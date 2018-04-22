package textRenderer

import (
	"go-experiments/common/opengl"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type TextRendererProgram struct {
	shaderProgram uint32
	projectionLoc int32
	cameraLoc     int32
	modelLoc      int32
	fontImageLoc  int32

	foregroundColorLoc int32
	backgroundColorLoc int32
}

func NewTextRendererProgram() *TextRendererProgram {
	var program TextRendererProgram
	program.shaderProgram = commonOpenGl.CreateProgram("./text/renderer/textRenderer")

	program.projectionLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("projection\x00"))
	program.cameraLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("camera\x00"))
	program.modelLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("model\x00"))
	program.fontImageLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("fontImage\x00"))

	program.foregroundColorLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("foregroundColor\x00"))
	program.backgroundColorLoc = gl.GetUniformLocation(program.shaderProgram, gl.Str("backgroundColor\x00"))

	return &program
}

func (r *TextRendererProgram) UseProgram(buffers *TextProgramBuffers) {
	gl.UseProgram(r.shaderProgram)
	gl.BindVertexArray(buffers.vao)
}

func (r *TextRendererProgram) SetColors(background, foreground mgl32.Vec3) {
	gl.Uniform3f(r.backgroundColorLoc, background.X(), background.Y(), background.Z())
	gl.Uniform3f(r.foregroundColorLoc, foreground.X(), foreground.Y(), foreground.Z())
}

func (s *TextRendererProgram) SetModel(model *mgl32.Mat4) {
	gl.UniformMatrix4fv(s.modelLoc, 1, false, &model[0])
}

func (r *TextRendererProgram) SetTexture(textureId, textureResource uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureResource)
	gl.Uniform1i(r.fontImageLoc, int32(textureId))
}

func (r *TextRendererProgram) Delete() {
	gl.DeleteProgram(r.shaderProgram)
}

// Implement Renderer
func (r *TextRendererProgram) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(r.shaderProgram)
	gl.UniformMatrix4fv(r.projectionLoc, 1, false, &projection[0])
}

func (r *TextRendererProgram) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(r.shaderProgram)
	gl.UniformMatrix4fv(r.cameraLoc, 1, false, &camera[0])
}
