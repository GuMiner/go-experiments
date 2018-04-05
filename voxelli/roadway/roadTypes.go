package roadway

import (
	"go-experiments/voxelli/geometry"
	"go-experiments/voxelli/utils"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Defines the default road type that is never in bounds
type OutOfBoundsRoad struct {
}

func (oob OutOfBoundsRoad) InBounds(position mgl32.Vec2) bool {
	return false
}

func (oob OutOfBoundsRoad) GetBounds(gridPos utils.IntVec2) []geometry.Intersectable {
	return []geometry.Intersectable{}
}

// Defines a road that is oriented in the X+ direction
type StraightRoad struct {
	rotated bool // If true, the road is oriented in the Y+ direction instead
}

func (straightRoad StraightRoad) InBounds(position mgl32.Vec2) bool {
	if straightRoad.rotated {
		// Limitation is in the Y direction
		return position.Y() > 1.0 && position.Y() < float32(GetGridSize()-1)
	}

	return position.X() > 1.0 && position.X() < float32(GetGridSize()-1)
}

func (straightRoad StraightRoad) GetBounds(gridPos utils.IntVec2) []geometry.Intersectable {
	gridSize := float32(GetGridSize())
	gridOffset := mgl32.Vec2{float32(gridPos.X()) * gridSize, float32(gridPos.Y()) * gridSize}

	if straightRoad.rotated {
		// Limitation on the top and the bottom
		return []geometry.Intersectable{
			geometry.NewLineSegment(mgl32.Vec2{0, gridSize}.Add(gridOffset), mgl32.Vec2{gridSize, gridSize}.Add(gridOffset)),
			geometry.NewLineSegment(mgl32.Vec2{0, 0}.Add(gridOffset), mgl32.Vec2{gridSize, 0}.Add(gridOffset)),
		}
	}

	return []geometry.Intersectable{
		geometry.NewLineSegment(mgl32.Vec2{gridSize, 0}.Add(gridOffset), mgl32.Vec2{gridSize, gridSize}.Add(gridOffset)),
		geometry.NewLineSegment(mgl32.Vec2{0, 0}.Add(gridOffset), mgl32.Vec2{0, gridSize}.Add(gridOffset)),
	}
}

// Defines a curved road that starts out straight in the X- direction and curves in the Y+ direction
type CurvedRoad struct {
	rotation int // Defines the number of times to rotate in the CW direction
}

func (curvedRoad CurvedRoad) InBounds(position mgl32.Vec2) bool {
	gridExtent := float32(GetGridSize()) - 0.5

	var length float32
	switch curvedRoad.rotation % 4 {
	case 0:
		// Center == (gridExtent, 0) (down-to-right)
		length = position.Sub(mgl32.Vec2{gridExtent, 0}).Len()
	case 1:
		// Center == (gridExtent, gridExtent) (up-to-right)
		length = position.Sub(mgl32.Vec2{gridExtent, gridExtent}).Len()
	case 2:
		// Center == (0, gridExtent) (left-to-up)
		length = position.Sub(mgl32.Vec2{0, gridExtent}).Len()
	case 3:
		// Center == (0, 0) (left-to-down)
		length = position.Len()
	}

	return length < gridExtent
}

func (curvedRoad CurvedRoad) GetBounds(gridPos utils.IntVec2) []geometry.Intersectable {
	gridSize := float32(GetGridSize())
	gridOffset := mgl32.Vec2{float32(gridPos.X()) * gridSize, float32(gridPos.Y()) * gridSize}

	var arc geometry.ArcSegment
	switch curvedRoad.rotation % 4 {
	case 0:
		arc = geometry.NewArcSegment(mgl32.Vec2{gridSize, 0}.Add(gridOffset), gridSize, math.Pi/2, math.Pi)
	case 1:
		arc = geometry.NewArcSegment(mgl32.Vec2{gridSize, gridSize}.Add(gridOffset), gridSize, math.Pi, 3*math.Pi/2)
	case 2:
		arc = geometry.NewArcSegment(mgl32.Vec2{0, gridSize}.Add(gridOffset), gridSize, 3*math.Pi/2, 2*math.Pi)
	case 3:
		arc = geometry.NewArcSegment(mgl32.Vec2{0, 0}.Add(gridOffset), gridSize, 0, math.Pi/2)
	}

	return []geometry.Intersectable{arc}
}
