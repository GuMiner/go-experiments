package main

import (
	"go-experiments/voxelli/voxel"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Defines how to render a voxelObject
type VoxelObjectRenderer struct {
	shaderProgram uint32

	projectionLoc    int32
	cameraLoc        int32
	modelLoc         int32
	timeLoc          int32
	colorOverrideLoc int32

	cubeArchetype *Cube
}

func (renderer *VoxelObjectRenderer) Render(voxelObject *voxel.VoxelObject, model *mgl32.Mat4) {
	for _, subObject := range voxelObject.SubObjects {
		for _, voxel := range subObject.Voxels {
			colorVector := voxelObject.Palette.Colors[voxel.ColorIdx-1].AsFloatVector()
			gl.Uniform4fv(renderer.colorOverrideLoc, 1, &colorVector[0])

			updatedModel := model.Mul4(mgl32.Translate3D(float32(voxel.Position.X()), float32(voxel.Position.Y()), float32(voxel.Position.Z())))
			gl.UniformMatrix4fv(renderer.modelLoc, 1, false, &updatedModel[0])
			renderer.cubeArchetype.Render()
		}
	}
}

func (renderer *VoxelObjectRenderer) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.projectionLoc, 1, false, &projection[0])
}

func (renderer *VoxelObjectRenderer) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.cameraLoc, 1, false, &camera[0])
}

func (renderer *VoxelObjectRenderer) UpdatedTime(time float32) {
	gl.UseProgram(renderer.shaderProgram)
	gl.Uniform1f(renderer.timeLoc, time)
}

func (renderer *VoxelObjectRenderer) Delete() {
	gl.DeleteProgram(renderer.shaderProgram)
	renderer.cubeArchetype.Delete()
}

func NewVoxelObjectRenderer() *VoxelObjectRenderer {
	var renderer VoxelObjectRenderer

	renderer.shaderProgram = createProgram("./shaders/basicRenderer")

	// Get locations of everything used in this program.
	renderer.projectionLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("projection\x00"))
	renderer.cameraLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("camera\x00"))
	renderer.modelLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("model\x00"))
	renderer.timeLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("runTime\x00"))
	renderer.colorOverrideLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("colorOverride\x00"))

	renderer.cubeArchetype = NewCube()

	return &renderer
}
