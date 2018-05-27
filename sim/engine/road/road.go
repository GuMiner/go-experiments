package road

import (
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

// TODO: Deduplicate with PowerLine, both are effectively equivalent
type RoadLine struct {
	start    mgl32.Vec2
	end      mgl32.Vec2
	capacity int64

	// The grid nodes this powerline refer to.
	startNode int
	endNode   int

	// Whether or not this road line is unconnected at the end, where
	//  it owns its endpoints.
	ownsStartNode bool
	ownsEndNode   bool
}

// Implement Element
// Gets the central position of the element.
func (p *RoadLine) GetRegion() *commonMath.Region {
	return nil
}

// Gets positions on the map that can be used to snap to points of the element.
func (p *RoadLine) GetSnapNodes() []mgl32.Vec2 {
	return []mgl32.Vec2{
		p.start,
		p.end}
}

// Gets lines on the map that can be used to snap to *edges* of the element
func (p *RoadLine) GetSnapEdges() [][2]mgl32.Vec2 {
	return [][2]mgl32.Vec2{
		[2]mgl32.Vec2{
			p.start,
			p.end}}
}

// Gets the line this road line represents.
func (p *RoadLine) GetLine() [2]mgl32.Vec2 {
	return p.GetSnapEdges()[0]
}

func (p *RoadLine) GetSnapNodeElement(snapNodeIdx int) int {
	if snapNodeIdx == 0 {
		return p.startNode
	}

	return p.endNode
}
