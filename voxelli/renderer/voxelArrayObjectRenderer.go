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

	projectionLoc int32
	cameraLoc     int32
	modelLoc      int32
}

func (renderer *VoxelArrayObjectRenderer) Render(voxelArrayObject *voxelArray.VoxelArrayObject, model *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)

	gl.UniformMatrix4fv(renderer.modelLoc, 1, false, &model[0])
	voxelArrayObject.Render()
}

func (renderer *VoxelArrayObjectRenderer) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.projectionLoc, 1, false, &projection[0])
}

func (renderer *VoxelArrayObjectRenderer) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.cameraLoc, 1, false, &camera[0])
}

func (renderer *VoxelArrayObjectRenderer) Delete() {
	gl.DeleteProgram(renderer.shaderProgram)
}

func NewVoxelArrayObjectRenderer() *VoxelArrayObjectRenderer {
	var renderer VoxelArrayObjectRenderer

	renderer.shaderProgram = opengl.CreateProgram("./shaders/voxelArrayRenderer")

	// Get locations of everything used in this program.
	renderer.projectionLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("projection\x00"))
	renderer.cameraLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("camera\x00"))
	renderer.modelLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("model\x00"))

	return &renderer
}
