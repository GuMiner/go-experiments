package power

import (
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

type PowerPlant struct {
	location mgl32.Vec2

	plantType string
	namedSize PowerPlantSize

	size        float32 // All plants are assumed square, for now. (TODO)
	orientation float32

	output int // kW
	// TODO: Add capacity factor, things impacted by, required resources, etc.

	gridId int
}

// Implement Element
// Gets the central position of the element
func (p *PowerPlant) GetRegion() *commonMath.Region {
	return &commonMath.Region{
		RegionType:  commonMath.SquareRegion,
		Position:    p.location,
		Scale:       p.size,
		Orientation: p.orientation}
}

// Gets positions on the map that can be used to snap to points of the element.
func (p *PowerPlant) GetSnapNodes() []mgl32.Vec2 {
	// TODO: This should be plant-type specific so it matches with the 2D model or image.
	nodes := make([]mgl32.Vec2, 1)
	nodes[0] = p.location
	return nodes
}

// Gets lines on the map that can be used to snap to *edges* of the element
func (p *PowerPlant) GetSnapEdges() [][2]mgl32.Vec2 {
	return make([][2]mgl32.Vec2, 0)
}

func (p *PowerPlant) GetSnapElement() int {
	return p.gridId
}
