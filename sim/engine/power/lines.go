package power

import (
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

type PowerLine struct {
	start    mgl32.Vec2
	end      mgl32.Vec2
	capacity int64

	// The grid nodes this powerline refer to.
	startNode int
	endNode   int

	// Whether or not this power line is unconnected at the end, where
	//  it owns its endpoints.
	ownsStartNode bool
	ownsEndNode   bool
}

// Implement Element
// Gets the central position of the element.
func (p *PowerLine) GetRegion() *commonMath.Region {
	return nil
}

// Gets positions on the map that can be used to snap to points of the element.
func (p *PowerLine) GetSnapNodes() []mgl32.Vec2 {
	return []mgl32.Vec2{
		p.start,
		p.end}
}

// Gets lines on the map that can be used to snap to *edges* of the element
func (p *PowerLine) GetSnapEdges() [][2]mgl32.Vec2 {
	return [][2]mgl32.Vec2{
		[2]mgl32.Vec2{
			p.start,
			p.end}}
}

// Gets the line this power line represents.
func (p *PowerLine) GetLine() [2]mgl32.Vec2 {
	return p.GetSnapEdges()[0]
}

func (p *PowerLine) GetSnapNodeElement(snapNodeIdx int) int {
	if snapNodeIdx == 0 {
		return p.startNode
	}

	return p.endNode
}
