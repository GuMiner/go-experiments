package region

import (
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

const circlePieces = 32

var ccwSquareRegion = []mgl32.Vec2{
	mgl32.Vec2{-0.5, 0.5},
	mgl32.Vec2{0.5, -0.5},
	mgl32.Vec2{0.5, 0.5},

	mgl32.Vec2{-0.5, -0.5},
	mgl32.Vec2{0.5, -0.5},
	mgl32.Vec2{-0.5, 0.5}}

func generateRegion(regionType commonMath.RegionType) []mgl32.Vec2 {
	switch regionType {
	case commonMath.SquareRegion:
		return ccwSquareRegion
	case commonMath.TriangleRegion:
		return generateTriangleRegion()
	default: // Circle region
		return generateCircleRegion()
	}
}

func generateTriangleRegion() []mgl32.Vec2 {
	vertices := make([]mgl32.Vec2, 3)
	vertices[0] = mgl32.Vec2{0.5, 0}
	vertices[1] = mgl32.Rotate2D(mgl32.DegToRad(120)).Mul2x1(vertices[0])
	vertices[2] = mgl32.Rotate2D(mgl32.DegToRad(240)).Mul2x1(vertices[0])
	return vertices
}

func generateCircleRegion() []mgl32.Vec2 {
	vertices := make([]mgl32.Vec2, circlePieces)
	for i := 0; i < circlePieces; i++ {
		rotation := mgl32.Rotate2D(mgl32.DegToRad(float32(i) * 360.0 / circlePieces))
		vertices[i] = rotation.Mul2x1(mgl32.Vec2{0.5, 0})
	}

	smallRotation := mgl32.Rotate2D(mgl32.DegToRad(360.0 / circlePieces))
	expandedVertices := make([]mgl32.Vec2, circlePieces*3)
	for i := 0; i < circlePieces; i++ {
		expandedVertices[i*3] = mgl32.Vec2{0, 0}
		expandedVertices[i*3+1] = vertices[i]
		expandedVertices[i*3+2] = smallRotation.Mul2x1(vertices[i])
	}

	return expandedVertices
}
