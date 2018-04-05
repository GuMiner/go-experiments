package main

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

func (oob OutOfBoundsRoad) FindBoundary(position mgl32.Vec2, direction mgl32.Vec2) (inBounds bool, gridId utils.IntVec2, gridPos mgl32.Vec2) {
	return true, utils.IntVec2{0, 0}, position
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

// Given a position and direction on the piece, finds the piece boundary.
// If the boundary is out-of-bounds, returns (true, {0, 0}, and the relative position of the intersection)
// If the boundary leads to another piece, returns (false, the offset to the next grid pos, and the relative position of the intersection on that next piece)
func (straightRoad StraightRoad) FindBoundary(position mgl32.Vec2, direction mgl32.Vec2) (inBounds bool, gridId utils.IntVec2, gridPos mgl32.Vec2) {
	gridSize := float32(GetGridSize())
	currentVector := geometry.NewVector(position, direction)

	// Determine which of the four walls of the grid this vector intersects.
	// We add in the direction (which is normalized) when we move into new pieces to prevent hitting on the wall we transferred through.
	left := geometry.NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{0, gridSize})
	if doesIntersect, intersectPos := left.Intersects(currentVector); doesIntersect {
		if straightRoad.rotated {
			// No limitation on the side, so we pass thru to the next piece, modifying the intersection position accordingly
			return false, utils.IntVec2{-1, 0}, intersectPos.Add(mgl32.Vec2{gridSize, 0}).Sub(direction)
		} else {
			return true, utils.IntVec2{0, 0}, intersectPos
		}
	}

	// Determine which of the four walls of the grid this vector intersects.
	right := geometry.NewLineSegment(mgl32.Vec2{gridSize, 0}, mgl32.Vec2{gridSize, gridSize})
	if doesIntersect, intersectPos := right.Intersects(currentVector); doesIntersect {
		if straightRoad.rotated {
			// No limitation on the side, so we pass thru to the next piece, modifying the intersection position accordingly
			return false, utils.IntVec2{1, 0}, intersectPos.Sub(mgl32.Vec2{gridSize, 0}).Add(direction)
		} else {
			return true, utils.IntVec2{0, 0}, intersectPos
		}
	}

	top := geometry.NewLineSegment(mgl32.Vec2{0, gridSize}, mgl32.Vec2{gridSize, gridSize})
	if doesIntersect, intersectPos := top.Intersects(currentVector); doesIntersect {
		if !straightRoad.rotated {
			// No limitation on the side, so we pass thru to the next piece, modifying the intersection position accordingly
			return false, utils.IntVec2{0, 1}, intersectPos.Sub(mgl32.Vec2{0, gridSize}).Add(direction)
		} else {
			return true, utils.IntVec2{0, 0}, intersectPos
		}
	}

	bottom := geometry.NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{gridSize, 0})
	if doesIntersect, intersectPos := bottom.Intersects(currentVector); doesIntersect {
		if !straightRoad.rotated {
			// No limitation on the side, so we pass thru to the next piece, modifying the intersection position accordingly
			return false, utils.IntVec2{0, -1}, intersectPos.Add(mgl32.Vec2{0, gridSize}).Sub(direction)
		} else {
			return true, utils.IntVec2{0, 0}, intersectPos
		}
	}

	return true, utils.IntVec2{0, 0}, position
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

// Given a position and direction on the piece, finds the piece boundary.
// If the boundary is out-of-bounds, returns (true, {0, 0}, and the relative position of the intersection)
// If the boundary leads to another piece, returns (false, the offset to the next grid pos, and the relative position of the intersection on that next piece)
func (curvedRoad CurvedRoad) FindBoundary(position mgl32.Vec2, direction mgl32.Vec2) (inBounds bool, gridId utils.IntVec2, gridPos mgl32.Vec2) {
	gridSize := float32(GetGridSize()) - 0.5
	currentVector := geometry.NewVector(position, direction)

	var arc geometry.ArcSegment
	switch curvedRoad.rotation % 4 {
	case 0:
		arc = geometry.NewArcSegment(mgl32.Vec2{gridSize, 0}, gridSize, math.Pi/2, math.Pi)
	case 1:
		arc = geometry.NewArcSegment(mgl32.Vec2{gridSize, gridSize}, gridSize, math.Pi, 3*math.Pi/2)
	case 2:
		arc = geometry.NewArcSegment(mgl32.Vec2{0, gridSize}, gridSize, 3*math.Pi/2, 2*math.Pi)
	case 3:
		arc = geometry.NewArcSegment(mgl32.Vec2{0, 0}, gridSize, 0, math.Pi/2)
	}

	if intersects, intersectionPoint := arc.Intersects(currentVector); intersects {
		// For each segment, we need to check two sides to see if the vector is heading to a new piece.
		return true, utils.IntVec2{0, 0}, intersectionPoint
	}

	return true, utils.IntVec2{0, 0}, position
}
