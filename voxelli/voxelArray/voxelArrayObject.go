package voxelArray

import (
	"go-experiments/voxelli/voxel"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// Defines a voxel object that has been condensed into a single VAO / VBO set
type VoxelArray struct {
	vao         uint32
	positionVbo uint32
	normalVbo   uint32
	colorVbo    uint32
}

func (v *VoxelArray) Delete() {
	gl.DeleteBuffers(1, &v.colorVbo)
	gl.DeleteBuffers(1, &v.normalVbo)
	gl.DeleteBuffers(1, &v.positionVbo)
	gl.DeleteVertexArrays(1, &v.vao)
}

func NewVoxelArray() *VoxelArray {
	var voxelArray VoxelArray
	gl.GenVertexArrays(1, &voxelArray.vao)
	gl.BindVertexArray(voxelArray.vao)

	gl.EnableVertexAttribArray(0)
	gl.GenBuffers(1, &voxelArray.positionVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, voxelArray.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.EnableVertexAttribArray(1)
	gl.GenBuffers(1, &voxelArray.normalVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, voxelArray.normalVbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.EnableVertexAttribArray(2)
	gl.GenBuffers(1, &voxelArray.colorVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, voxelArray.colorVbo)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)

	return &voxelArray
}

type VoxelArrayObject struct {
	voxelArray  *VoxelArray
	VoxelObject *voxel.VoxelObject

	Vertices int
}

// Collapses a series of voxels into a voxel array for speedy rendering
func (v *VoxelArrayObject) condense() {
	collapsedVoxels := collapseVoxels(v.VoxelObject)

	v.Vertices = len(collapsedVoxels.positionVertices)

	// *3 == vec3, *4 == float32
	gl.BindBuffer(gl.ARRAY_BUFFER, v.voxelArray.positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, v.Vertices*3*4, gl.Ptr(collapsedVoxels.positionVertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, v.voxelArray.normalVbo)
	gl.BufferData(gl.ARRAY_BUFFER, v.Vertices*3*4, gl.Ptr(collapsedVoxels.normalVertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, v.voxelArray.colorVbo)
	gl.BufferData(gl.ARRAY_BUFFER, v.Vertices*3*4, gl.Ptr(collapsedVoxels.colorVertices), gl.STATIC_DRAW)
}

func (v *VoxelArrayObject) Render() {
	gl.BindVertexArray(v.voxelArray.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(v.Vertices))
}

func (v *VoxelArrayObject) Delete() {
	v.voxelArray.Delete()
}

func NewVoxelArrayObject(voxelObject *voxel.VoxelObject) *VoxelArrayObject {
	voxelArrayObject := VoxelArrayObject{VoxelObject: voxelObject}
	voxelArrayObject.voxelArray = NewVoxelArray()
	voxelArrayObject.condense()

	return &voxelArrayObject
}
