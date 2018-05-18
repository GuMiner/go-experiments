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

// Iterates through all integer positions in the region, calling iterate.
// iterate() returns true to exit early.
// This function returns true if exited early, false otherwise.
func (r Region) IterateIntWithEarlyExit(iterate func(x, y int) bool) bool {

	switch r.RegionType {
	case SquareRegion:
		// TODO: Account for orientation.
		for i := int(r.Position.X() - r.Scale/2); i <= int(r.Position.X()+r.Scale/2); i++ {
			for j := int(r.Position.Y() - r.Scale/2); j <= int(r.Position.Y()+r.Scale/2); j++ {
				if iterate(i, j) {
					return true
				}
			}
		}
	case CircleRegion:
		// TODO: Make this more efficient
		for i := int(r.Position.X() - r.Scale/2); i <= int(r.Position.X()+r.Scale/2); i++ {
			for j := int(r.Position.Y() - r.Scale/2); j <= int(r.Position.Y()+r.Scale/2); j++ {
				if ((float32(i)-r.Position.X())*(float32(i)-r.Position.X()) + (float32(j)-r.Position.Y())*(float32(j)-r.Position.Y())) < r.Scale*r.Scale/4 {
					if iterate(i, j) {
						return true
					}
				}
			}
		}
		return false
	default: // Triangle Region
		// TODO:
		return false
	}

	return false
}
