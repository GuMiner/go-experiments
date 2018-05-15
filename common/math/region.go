package commonMath

import (
	"github.com/go-gl/mathgl/mgl32"
)

type RegionType int

const (
	CircleRegion RegionType = iota
	SquareRegion
	TriangleRegion
)

// Defines a region
type Region struct {
	RegionType RegionType

	Scale       float32
	Position    mgl32.Vec2
	Orientation float32
}
