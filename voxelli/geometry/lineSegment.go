package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type LineSegment struct {
	start mgl32.Vec2
	end   mgl32.Vec2
}

// Returns the % multiplier to intersect both vectors. Vectors do NOT need to be normalized.
func intersectionImplementation(first, second Vector) (float32, float32) {
	// first.pos + r*first.direction = second.pos + s*seciond.direction
	slope := second.direction.Y() / second.direction.X()
	r := ((first.point.Y() - second.point.Y()) + slope*(second.point.X()-first.point.X())) / (first.direction.X()*slope - first.direction.Y())
	s := (first.point.X() + r*first.direction.X() - second.point.X()) / second.direction.X()

	return r, s
}

// Returns true and the intersection point on an intersection, false otherwise
func (seg LineSegment) Intersects(vector Vector) (bool, mgl32.Vec2, mgl32.Vec2) {
	// Convert the line segment to vector form.
	lineVector := Vector{point: seg.start, direction: seg.end.Sub(seg.start)}

	if math.Abs(float64(lineVector.direction.Normalize().Dot(vector.direction))) > 0.9999 { // Small wiggle factor for floating point precision
		// The line and the vector are essentially pointing in the same direction, so they can be coincident but don't intersect
		return false, mgl32.Vec2{0, 0}, mgl32.Vec2{0, 0}
	}

	var r, s float32 // r == % distance on the line. s == % distance on the vector
	if math.Abs(float64(vector.direction.X())) < math.Abs(float64(lineVector.direction.X())) {

		// Swap the vector and the line vector to avoid dividing by zero
		s, r = intersectionImplementation(vector, lineVector)
	} else {
		r, s = intersectionImplementation(lineVector, vector)
	}

	// We are not on the line or we intersected backwards from the vector
	if s < 0 || r < 0 || r > 1 {
		return false, mgl32.Vec2{0, 0}, mgl32.Vec2{0, 0}
	}

	return true, lineVector.point.Add(lineVector.direction.Mul(r)), lineVector.direction.Normalize()
}

func NewLineSegment(start, end mgl32.Vec2) LineSegment {
	return LineSegment{start: start, end: end}
}
