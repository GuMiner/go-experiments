package element

import "github.com/go-gl/mathgl/mgl32"

// Defines commonalities all interactable map elements have.
type Element interface {
	// Gets the type of map element
	GetType() MapElementType

	// Gets the central position of the element
	GetPosition() mgl32.Vec2

	// Gets positions on the map that can be used to points of the element.
	GetSnapNodes() []mgl32.Vec2

	// Gets lines on the map that can be used to snap to *edges* of the element
	GetSnapEdges() [][2]mgl32.Vec2
}
