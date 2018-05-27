package road

import (
	"fmt"
	"go-experiments/sim/engine/core"
	"go-experiments/sim/engine/element"

	"github.com/go-gl/mathgl/mgl32"
)

type RoadGrid struct {
	grid    *core.ResizeableGraph
	nodeMap map[int]element.Element // Reverse maps a node ID to an element

	roadLines    map[int64]*RoadLine
	nextRoadLine int64
}

func NewRoadGrid() *RoadGrid {
	grid := RoadGrid{
		grid:         core.NewResizeableGraph(),
		nodeMap:      make(map[int]element.Element),
		roadLines:    make(map[int64]*RoadLine),
		nextRoadLine: 0}

	return &grid
}

func (p *RoadGrid) AddLine(start, end mgl32.Vec2, capacity int64, startNode, endNode int) *RoadLine {
	line := RoadLine{
		start:    start,
		end:      end,
		capacity: capacity}

	if startNode == endNode && startNode != -1 {
		fmt.Printf("Roads must be between nodes and cannot (for a single line) loop\n")
		return nil
	} else if startNode != -1 && endNode != -1 {
		// This might be a duplicate line.
		cost := p.grid.Cost(startNode, endNode)
		if cost != -1 {
			fmt.Printf("There already is a road from %v to %v.\n", startNode, endNode)
			return nil
		}
	}

	if startNode == -1 {
		line.startNode = p.grid.AddNode()
		line.ownsStartNode = true
		p.nodeMap[line.startNode] = &line
	} else {
		line.startNode = startNode
		line.ownsStartNode = false
	}

	if endNode == -1 {
		line.endNode = p.grid.AddNode()
		line.ownsEndNode = true
		p.nodeMap[line.endNode] = &line
	} else {
		line.endNode = endNode
		line.ownsEndNode = false
	}

	p.grid.AddOrUpdateEdgeCost(line.startNode, line.endNode, line.capacity)
	p.grid.AddOrUpdateEdgeCost(line.endNode, line.startNode, line.capacity)
	p.roadLines[p.nextRoadLine] = &line
	p.nextRoadLine++

	return &line
}

func (p *RoadGrid) IterateLines(iterate func(*RoadLine)) {
	for _, line := range p.roadLines {
		iterate(line)
	}
}
