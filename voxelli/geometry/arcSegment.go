package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type ArcSegment struct {
	center     mgl32.Vec2
	radius     float32
	angleStart float32
	angleEnd   float32
}

// Solves the quadratic equation, assuming that a != 0
// Returns the listing of real results (0 to 2)
func solveQuadraticReals(a, b, c float32) []float32 {
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []float32{}
	} else if discriminant == 0 {
		return []float32{-b / (2 * a)}
	}

	discriminantSqrt := float32(math.Sqrt(float64(discriminant)))
	return []float32{
		(-b + discriminantSqrt) / (2 * a),
		(-b - discriminantSqrt) / (2 * a)}
}

// Finds the closest vector intersection to the given arc segment as a circle.
// This can return up to two results, so we take the smallest positive result for the
// % along the vector direction if we get any results.
// See the associated wxMaxima file for the mathematical basis for this.
func findVectorCircleIntersection(seg *ArcSegment, vector Vector) []float32 {
	xv := vector.direction.X()
	yv := vector.direction.Y()

	xp := vector.point.X()
	yp := vector.point.Y()
	xo := seg.center.X()
	yo := seg.center.Y()

	a := xv*xv + yv*yv
	b := 2 * (xp*xv + yp*yv - (xv*xo + yv*yo))
	c := xp*xp + xo*xo + yp*yp + yo*yo - (2*(xp*xo+yp*yo) + seg.radius*seg.radius)

	// Find the intersection point and reduce to the closest positive point
	// TODO Tomorrow -- this isn't correct, we need to return both results.
	// If we don't return both results, we match on a closer part of the circle that may not be part of the partial arc.
	// This should actually simplify the logic a bit.
	s := solveQuadraticReals(a, b, c)
	return s
}

// Returns true and the intersection point on an intersection, false otherwise
func (seg *ArcSegment) Intersects(vector Vector) (bool, mgl32.Vec2) {
	intersectionPercentDistances := findVectorCircleIntersection(seg, vector)

	doesIntersect := false
	foundIntersectionPoint := mgl32.Vec2{0, 0}
	lastIntersectionDistancePercent := float32(math.MaxFloat32)
	for _, intersectionDistancePercent := range intersectionPercentDistances {
		if intersectionDistancePercent > 0 && intersectionDistancePercent < lastIntersectionDistancePercent {
			intersectionPoint := mgl32.Vec2{
				vector.point.X() + vector.direction.X()*intersectionDistancePercent,
				vector.point.Y() + vector.direction.Y()*intersectionDistancePercent}

			// Verify the intersection point is on the arc.
			dx := intersectionPoint.X() - seg.center.X()
			dy := intersectionPoint.Y() - seg.center.Y()

			angle := float32(math.Atan2(float64(dy), float64(dx)))
			if angle < 0 {
				angle += 2 * math.Pi
			}

			if angle >= seg.angleStart && angle <= seg.angleEnd {
				doesIntersect = true
				foundIntersectionPoint = intersectionPoint
				lastIntersectionDistancePercent = intersectionDistancePercent
			}
		}
	}

	return doesIntersect, foundIntersectionPoint
}

func NewArcSegment(center mgl32.Vec2, radius, angleStart, angleEnd float32) ArcSegment {
	return ArcSegment{center: center, radius: radius, angleStart: angleStart, angleEnd: angleEnd}
}
