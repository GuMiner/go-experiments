package element

import (
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

// Defines commonalities all interactable map elements have.
type Element interface {
	// Gets a region describing the position / orientation / size / type of location this is.
	GetRegion() *commonMath.Region

	// Gets positions on the map that can be used to snap to points of the element.
	GetSnapNodes() []mgl32.Vec2

	// Gets lines on the map that can be used to snap to *edges* of the element
	GetSnapEdges() [][2]mgl32.Vec2
}
