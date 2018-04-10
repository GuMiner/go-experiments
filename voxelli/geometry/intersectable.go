package geometry

import "github.com/go-gl/mathgl/mgl32"

type Intersectable interface {
	// Given a vector, returns true if it intersects, the point of intersection, and normal of intersection
	Intersects(vector Vector) (bool, mgl32.Vec2, mgl32.Vec2)
}
