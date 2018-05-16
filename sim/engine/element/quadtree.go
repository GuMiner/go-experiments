package element

import (
	"fmt"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/subtile"

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

type PointWithDistance struct {
	distance float32
	point    mgl32.Vec2
	metadata interface{}
}

func NewPointWithDistance(point qtree.Point, referencePoint mgl32.Vec2) PointWithDistance {
	// TODO: See below, this won't work.
	return PointWithDistance{}
}

func NewMultiQuadtree() *MultiQuadtree {
	return &MultiQuadtree{
		quadtrees: make(map[int]map[int]*qtree.Qtree)}
}

func (m *MultiQuadtree) Add(elem Element) {
	region := elem.GetRegion()

	regionX, regionY := subtile.GetRegionIndices(region.Position, config.Config.Terrain.RegionSize)

	if _, ok := m.quadtrees[regionX]; !ok {
		m.quadtrees[regionX] = make(map[int]*qtree.Qtree)
	}

	if _, ok := m.quadtrees[regionX][regionY]; !ok {
		// We offset the quadtree so they all slightly overlap, to avoid edge effects.
		m.quadtrees[regionX][regionY] = qtree.NewQtree(
			-1,
			-1,
			float64(config.Config.Terrain.RegionSize+2),
			float64(config.Config.Terrain.RegionSize+2),
			maxPointsInBucket)
	}

	pos := subtile.GetLocalFloatIndices(region.Position, regionX, regionY, config.Config.Terrain.RegionSize)
	point := qtree.NewPoint(float64(pos.X()), float64(pos.Y()), elem)
	m.quadtrees[regionX][regionY].Insert(point)
}

// Returns the K-nearest points.
// If not enough points are found in the local quadtree region, searches surrounding reginons.
// Stops when it has searched the central and all surrounding regions (9 in total) or found K points.
func (m *MultiQuadtree) KNearest(pos mgl32.Vec2, count int) []Element {
	regionX, regionY := subtile.GetRegionIndices(pos, config.Config.Terrain.RegionSize)

	points := make([]PointWithDistance, 0)
	for i := regionX - 1; i <= regionX+1; i++ {
		for j := regionY - 1; j < regionY+1; j++ {
			if quadtree, ok := m.quadtrees[regionX][regionY]; ok {

				// TODO This won't work as the point by definition will only be in-bounds for one location
				// Basically, I should probably implement my own quadtree instead of attempting to push each framework into
				//  the approaches I need.
				mglPoint := subtile.GetLocalFloatIndices(pos, i, j, config.Config.Terrain.RegionSize)

				point := qtree.NewPoint(float64(mglPoint.X()), float64(mglPoint.Y()), nil)
				for _, point := range quadtree.KNN(point, count) {
					// points = append(points, NewPointWithDistance(point, mquadtree.KNN(point, count)...))
					fmt.Printf("%v", point)
				}
			}
		}
	}

	// TODO: Sort points by distance and restrict to count!
	elements := make([]Element, len(points))
	for i, point := range points {
		element, _ := point.metadata.(Element)
		elements[i] = element
	}

	return elements
}
