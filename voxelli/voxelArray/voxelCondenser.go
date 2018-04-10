package voxelArray

import (
	"go-experiments/voxelli/utils"
	"go-experiments/voxelli/voxel"

	"github.com/go-gl/mathgl/mgl32"
)

var XNeg = []mgl32.Vec3{
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5}}

var XPos = []mgl32.Vec3{
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, 0.5}}

var YNeg = []mgl32.Vec3{
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5}}

var YPos = []mgl32.Vec3{
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5}}

var ZNeg = []mgl32.Vec3{
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5}}

var ZPos = []mgl32.Vec3{
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5}}

type collapsedVertices struct {
	positionVertices []mgl32.Vec3
	normalVertices   []mgl32.Vec3
	colorVertices    []mgl32.Vec3
}

func (c *collapsedVertices) addPlane(positionVertices []mgl32.Vec3, offset utils.IntVec3, normal mgl32.Vec3, color mgl32.Vec3) {
	c.positionVertices = append(c.positionVertices, positionVertices...)
	for i := 0; i < len(positionVertices); i++ {
		j := (len(c.positionVertices) - len(positionVertices)) + i

		c.positionVertices[j] = c.positionVertices[j].Add(offset.AsFloatVector())
		c.normalVertices = append(c.normalVertices, normal)
		c.colorVertices = append(c.colorVertices, color)
	}
}

// Collapses a series of voxels into a voxel array for speedy rendering
func collapseVoxels(voxelObject *voxel.VoxelObject) *collapsedVertices {
	// Build a listing of what voxels exist where
	var voxelExistenceMap map[utils.IntVec3]bool = make(map[utils.IntVec3]bool)

	for _, subObject := range voxelObject.SubObjects {
		for _, voxel := range subObject.Voxels {
			voxelExistenceMap[voxel.Position] = true
		}
	}

	var vertices collapsedVertices
	for _, subObject := range voxelObject.SubObjects {
		for _, voxel := range subObject.Voxels {
			color := voxelObject.Palette.Colors[voxel.ColorIdx-1].AsOpaqueFloatVector()

			// Check each plane to see if there is a corresponding voxel. If there is, we don't add a plane for it as that plane is hidden right now.
			if !voxelExistenceMap[utils.IntVec3{voxel.Position.X() - 1, voxel.Position.Y(), voxel.Position.Z()}] {
				vertices.addPlane(XNeg, voxel.Position, mgl32.Vec3{-1, 0, 0}, color)
			}
			if !voxelExistenceMap[utils.IntVec3{voxel.Position.X() + 1, voxel.Position.Y(), voxel.Position.Z()}] {
				vertices.addPlane(XPos, voxel.Position, mgl32.Vec3{1, 0, 0}, color)
			}

			if !voxelExistenceMap[utils.IntVec3{voxel.Position.X(), voxel.Position.Y() - 1, voxel.Position.Z()}] {
				vertices.addPlane(YNeg, voxel.Position, mgl32.Vec3{0, -1, 0}, color)
			}
			if !voxelExistenceMap[utils.IntVec3{voxel.Position.X(), voxel.Position.Y() + 1, voxel.Position.Z()}] {
				vertices.addPlane(YPos, voxel.Position, mgl32.Vec3{0, 1, 0}, color)
			}

			if !voxelExistenceMap[utils.IntVec3{voxel.Position.X(), voxel.Position.Y(), voxel.Position.Z() - 1}] {
				vertices.addPlane(ZNeg, voxel.Position, mgl32.Vec3{0, 0, -1}, color)
			}
			if !voxelExistenceMap[utils.IntVec3{voxel.Position.X(), voxel.Position.Y(), voxel.Position.Z() + 1}] {
				vertices.addPlane(ZPos, voxel.Position, mgl32.Vec3{0, 0, 1}, color)
			}
		}
	}

	return &vertices
}
