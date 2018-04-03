package geometry

import "github.com/go-gl/mathgl/mgl32"

type Vector struct {
	point     mgl32.Vec2
	direction mgl32.Vec2
}

func NewVector(point, direction mgl32.Vec2) Vector {
	return Vector{point: point, direction: direction.Normalize()}
}
