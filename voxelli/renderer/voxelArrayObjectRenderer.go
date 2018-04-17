package renderer

import (
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/voxelArray"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Defines how to render a voxelObject
type VoxelArrayObjectRenderer struct {
	shaderProgram uint32

	projectionLoc   int32
	cameraLoc       int32
	modelLoc        int32
	shadingColorLoc int32

	partialShadowMatrixLoc int32
	shadowTextureLoc       int32
	shadowTextureId        uint32
}

func (renderer *VoxelArrayObjectRenderer) Render(voxelArrayObject *voxelArray.VoxelArrayObject, model *mgl32.Mat4, shadingColor mgl32.Vec3) {
	gl.UseProgram(renderer.shaderProgram)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, renderer.shadowTextureId)
	gl.Uniform1i(renderer.shadowTextureLoc, 0)

	gl.UniformMatrix4fv(renderer.modelLoc, 1, false, &model[0])
	gl.Uniform3f(renderer.shadingColorLoc, shadingColor[0], shadingColor[1], shadingColor[2])
	voxelArrayObject.Render()
}

func (renderer *VoxelArrayObjectRenderer) UpdateShadows(partialShadowMatrix *mgl32.Mat4, shadowTextureId uint32) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.partialShadowMatrixLoc, 1, false, &partialShadowMatrix[0])
	renderer.shadowTextureId = shadowTextureId
}

func (renderer *VoxelArrayObjectRenderer) Delete() {
	gl.DeleteProgram(renderer.shaderProgram)
}

func NewVoxelArrayObjectRenderer() *VoxelArrayObjectRenderer {
	var renderer VoxelArrayObjectRenderer

	renderer.shaderProgram = opengl.CreateProgram("./voxelArray/voxelArrayRenderer")

	// Get locations of everything used in this program.
	renderer.projectionLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("projection\x00"))
	renderer.cameraLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("camera\x00"))
	renderer.modelLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("model\x00"))
	renderer.shadingColorLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("shadingColor\x00"))

	renderer.partialShadowMatrixLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("partialShadowMatrix\x00"))
	renderer.shadowTextureLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("shadowTexture\x00"))

	return &renderer
}

// Implement Renderer
func (renderer *VoxelArrayObjectRenderer) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.projectionLoc, 1, false, &projection[0])
}

func (renderer *VoxelArrayObjectRenderer) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.cameraLoc, 1, false, &camera[0])
}
