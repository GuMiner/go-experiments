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
		for j := 0; j < len(roadway.roadElements[0]); j++ {
			switch val := roadway.roadElements[i][j].(type) {
			case StraightRoad:
				modelMatrix := mgl32.Translate3D(float32(GetGridSize()/2), float32(GetGridSize()/2), 0.0)

				if val.rotated {
					modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(90.0), mgl32.Vec3{0.0, 0.0, 1.0}))
				}

				modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(float32(i*GetGridSize()), float32(j*GetGridSize()), 0.0))
				renderer.voxelRenderer.Render(renderer.straightRoad, &modelMatrix)
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
