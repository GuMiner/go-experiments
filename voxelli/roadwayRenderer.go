package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// Defines how to render a roadway
type RoadwayRenderer struct {
	voxelRenderer *VoxelObjectRenderer
	straightRoad  *VoxelObject
	curvedRoad    *VoxelObject
}

func (renderer *RoadwayRenderer) Render(roadway *Roadway) {
	for i := 0; i < len(roadway.roadElements); i++ {
		for j := 0; j < len(roadway.roadElements[i]); j++ {
			positionMatrix := mgl32.Translate3D(float32(i*GetGridSize()), float32(j*GetGridSize()), 0.0)

			// Center the road piece at the start.
			rotateOffsetMatrix := mgl32.Translate3D(-float32(GetGridSize()/2)+0.5, -float32(GetGridSize()/2)+0.5, 0.0)

			switch val := roadway.roadElements[i][j].(type) {
			case StraightRoad:
				rotateMatrix := mgl32.Ident4()
				if val.rotated {
					rotateMatrix = mgl32.HomogRotate3D(mgl32.DegToRad(90.0), mgl32.Vec3{0.0, 0.0, 1.0})
				}

				modelMatrix := positionMatrix.Mul4(rotateMatrix.Mul4(rotateOffsetMatrix))
				renderer.voxelRenderer.Render(renderer.straightRoad, &modelMatrix)
			case CurvedRoad:
				rotateMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(90.0*float32(val.rotation)), mgl32.Vec3{0.0, 0.0, 1.0})

				modelMatrix := positionMatrix.Mul4(rotateMatrix.Mul4(rotateOffsetMatrix))
				renderer.voxelRenderer.Render(renderer.curvedRoad, &modelMatrix)
			}
		}
	}
}

func (renderer *RoadwayRenderer) loadRoadTypes() {
	renderer.straightRoad = NewVoxelObject("./data/models/road_straight.vox")
	fmt.Printf("Straight Road objects: %v\n", len(renderer.straightRoad.subObjects))

	renderer.curvedRoad = NewVoxelObject("./data/models/road_curved.vox")
	fmt.Printf("Curved Road objects: %v\n", len(renderer.curvedRoad.subObjects))
}

func NewRoadwayRenderer(voxelRenderer *VoxelObjectRenderer) *RoadwayRenderer {
	var renderer RoadwayRenderer
	renderer.voxelRenderer = voxelRenderer
	renderer.loadRoadTypes()

	return &renderer
}
