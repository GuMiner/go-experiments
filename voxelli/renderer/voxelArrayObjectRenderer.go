package renderer

import (
	"go-experiments/common/opengl"
	"go-experiments/voxelli/voxelArray"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type shaderProgramBase struct {
	program uint32

	projectionLoc int32
	cameraLoc     int32
	modelLoc      int32
}

// Defines how to render a voxelObject
type VoxelArrayObjectRenderer struct {
	shaderProgram   shaderProgramBase
	shadingColorLoc int32

	partialShadowMatrixLoc int32
	shadowTextureLoc       int32
	shadowTextureId        uint32

	depthShaderProgram shaderProgramBase
}

var depthModeOnly bool = false

func EnableDepthModeOnly() {
	depthModeOnly = true
}

func DisableDepthModeOnly() {
	depthModeOnly = false
}

func (renderer *VoxelArrayObjectRenderer) Render(voxelArrayObject *voxelArray.VoxelArrayObject, model *mgl32.Mat4, shadingColor mgl32.Vec3) {
	if depthModeOnly {
		gl.UseProgram(renderer.depthShaderProgram.program)
		gl.UniformMatrix4fv(renderer.depthShaderProgram.modelLoc, 1, false, &model[0])
	} else {
		gl.UseProgram(renderer.shaderProgram.program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, renderer.shadowTextureId)
		gl.Uniform1i(renderer.shadowTextureLoc, 0)

		gl.UniformMatrix4fv(renderer.shaderProgram.modelLoc, 1, false, &model[0])
		gl.Uniform3f(renderer.shadingColorLoc, shadingColor[0], shadingColor[1], shadingColor[2])
	}

	voxelArrayObject.Render()
}

func (renderer *VoxelArrayObjectRenderer) UpdateShadows(partialShadowMatrix *mgl32.Mat4, shadowTextureId uint32) {
	gl.UseProgram(renderer.shaderProgram.program)
	gl.UniformMatrix4fv(renderer.partialShadowMatrixLoc, 1, false, &partialShadowMatrix[0])

	renderer.shadowTextureId = shadowTextureId
}

func (renderer *VoxelArrayObjectRenderer) Delete() {
	gl.DeleteProgram(renderer.shaderProgram.program)
	gl.DeleteProgram(renderer.depthShaderProgram.program)
}

func NewVoxelArrayObjectRenderer() *VoxelArrayObjectRenderer {
	var renderer VoxelArrayObjectRenderer

	renderer.shaderProgram = shaderProgramBase{}
	renderer.shaderProgram.program = commonOpenGl.CreateProgram("./voxelArray/voxelArrayRenderer")

	renderer.shaderProgram.projectionLoc = gl.GetUniformLocation(renderer.shaderProgram.program, gl.Str("projection\x00"))
	renderer.shaderProgram.cameraLoc = gl.GetUniformLocation(renderer.shaderProgram.program, gl.Str("camera\x00"))
	renderer.shaderProgram.modelLoc = gl.GetUniformLocation(renderer.shaderProgram.program, gl.Str("model\x00"))
	renderer.shadingColorLoc = gl.GetUniformLocation(renderer.shaderProgram.program, gl.Str("shadingColor\x00"))

	renderer.partialShadowMatrixLoc = gl.GetUniformLocation(renderer.shaderProgram.program, gl.Str("partialShadowMatrix\x00"))
	renderer.shadowTextureLoc = gl.GetUniformLocation(renderer.shaderProgram.program, gl.Str("shadowTexture\x00"))

	renderer.depthShaderProgram = shaderProgramBase{}
	renderer.depthShaderProgram.program = commonOpenGl.CreateProgram("./voxelArray/voxelArrayDepthRenderer")
	renderer.depthShaderProgram.projectionLoc = gl.GetUniformLocation(renderer.depthShaderProgram.program, gl.Str("projection\x00"))
	renderer.depthShaderProgram.cameraLoc = gl.GetUniformLocation(renderer.depthShaderProgram.program, gl.Str("camera\x00"))
	renderer.depthShaderProgram.modelLoc = gl.GetUniformLocation(renderer.depthShaderProgram.program, gl.Str("model\x00"))

	return &renderer
}

// Implement Renderer
func (renderer *VoxelArrayObjectRenderer) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram.program)
	gl.UniformMatrix4fv(renderer.shaderProgram.projectionLoc, 1, false, &projection[0])

	gl.UseProgram(renderer.depthShaderProgram.program)
	gl.UniformMatrix4fv(renderer.depthShaderProgram.projectionLoc, 1, false, &projection[0])
}

func (renderer *VoxelArrayObjectRenderer) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram.program)
	gl.UniformMatrix4fv(renderer.shaderProgram.cameraLoc, 1, false, &camera[0])

	gl.UseProgram(renderer.depthShaderProgram.program)
	gl.UniformMatrix4fv(renderer.depthShaderProgram.cameraLoc, 1, false, &camera[0])
}
