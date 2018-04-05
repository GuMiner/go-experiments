package geometry

import "github.com/go-gl/mathgl/mgl32"

type Intersectable interface {
	Intersects(vector Vector) (bool, mgl32.Vec2)
}
