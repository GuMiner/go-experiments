package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
)

type ArcSegment struct {
	center     mgl32.Vec2
	radius     float32
	angleStart float32
	angleEnd   float32
}

// Returns true and the intersection point on an intersection, false otherwise
func (seg *ArcSegment) Intersects(vector Vector) (bool, mgl32.Vec2) {

	// TODO: implement
	return false, mgl32.Vec2{0, 0}
}

func NewArcSegment(center mgl32.Vec2, radius, angleStart, angleEnd float32) ArcSegment {
	return ArcSegment{center: center, radius: radius, angleStart: angleStart, angleEnd: angleEnd}
}
