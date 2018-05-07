package element

import (
	"go-experiments/sim/config"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/mleveck/go-quad-tree"
)

// Defines a quadtree for efficient area-based searching of elements
// This simulator doesn't have bounds, so we define a series of quadtrees as multiples of our regions.

// Empirically determined
const maxPointsInBucket = 16

type MultiQuadtree struct {
	quadtrees map[int]map[int]*qtree.Qtree
}

func NewMultiQuadtree() *MultiQuadtree {
	return &MultiQuadtree{
		quadtrees: make(map[int]map[int]*qtree.Qtree)}
}

// TODO: We need something so that when we search we can go back to stuff, like powerlines.
func (m *MultiQuadtree) AddPoint(pos mgl32.Vec2) *qtree.Point {
	quadtreeSize := config.Config.Terrain.RegionSize * 4
	subquadTreePos := pos.Mul(1.0 / float32(quadtreeSize))

	x := int(subquadTreePos.X())
	y := int(subquadTreePos.Y())
	if _, ok := m.quadtrees[x]; !ok {
		m.quadtrees[x] = make(map[int]*qtree.Qtree)
	}

	if _, ok := m.quadtrees[x][y]; !ok {
		m.quadtrees[x][y] = qtree.NewQtree(
			float64(x*quadtreeSize),
			float64(y*quadtreeSize),
			float64(quadtreeSize),
			float64(quadtreeSize),
			maxPointsInBucket)
	}

	point := qtree.NewPoint(float64(pos.X()), float64(pos.Y()), nil)
	m.quadtrees[x][y].Insert(point)

	return point
}
