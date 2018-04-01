package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Defines the default road type that is never in bounds
type OutOfBoundsRoad struct {
}

func (oob OutOfBoundsRoad) InBounds(position mgl32.Vec2) bool {
	return false
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

// Defines a curved road that starts out straight in the X- direction and curves in the Y+ direction

type CurvedRoad struct {
	rotation int // Defines the number of times to rotate in the CW direction
}

func (curvedRoad CurvedRoad) InBounds(position mgl32.Vec2) bool {
	gridExtent := float32(GetGridSize()) - 0.5

	var length float32
	switch curvedRoad.rotation % 4 {
	case 0:
		length = position.Sub(mgl32.Vec2{gridExtent, 0}).Len()
	case 1:
		length = position.Sub(mgl32.Vec2{gridExtent, gridExtent}).Len()
	case 2:
		length = position.Sub(mgl32.Vec2{0, gridExtent}).Len()
	case 3:
		length = position.Len()
	}

	return length < gridExtent
}
