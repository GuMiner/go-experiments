package power

import (
	"fmt"
	"go-experiments/sim/engine/core"
	"go-experiments/sim/engine/element"

	"github.com/go-gl/mathgl/mgl32"
)

type PowerGrid struct {
	grid    *core.ResizeableGraph
	nodeMap map[int]element.Element // Reverse maps a node ID to an element.

	powerPlants    map[int64]*PowerPlant
	nextPowerPlant int64

	powerLines    map[int64]*PowerLine
	nextPowerLine int64
}

func NewPowerGrid() *PowerGrid {
	grid := PowerGrid{
		grid:           core.NewResizeableGraph(),
		nodeMap:        make(map[int]element.Element),
		powerPlants:    make(map[int64]*PowerPlant),
		nextPowerPlant: 0,
		powerLines:     make(map[int64]*PowerLine),
		nextPowerLine:  0}

	return &grid
}

func (p *PowerGrid) Add(pos mgl32.Vec2, plantType string, plantSize PowerPlantSize) *PowerPlant {
	output, size := GetPowerOutputAndSize(plantType, plantSize)

	plant := PowerPlant{
		location:    pos,
		plantType:   plantType,
		namedSize:   plantSize,
		size:        float32(size),
		orientation: 0, // TODO: Rotation
		output:      output,
		gridId:      p.grid.AddNode()}

	p.nodeMap[plant.gridId] = &plant

	p.powerPlants[p.nextPowerPlant] = &plant
	fmt.Printf("Added power plant '%v'.\n", p.powerPlants[p.nextPowerPlant])

	p.nextPowerPlant++

	return &plant
}

// Adds a powerline. For both startNode and endNode, if -1 generates a new grid node, else uses an existing node.
func (p *PowerGrid) AddLine(start, end mgl32.Vec2, capacity int64, startNode, endNode int) *PowerLine {
	line := PowerLine{
		start:    start,
		end:      end,
		capacity: capacity}

	if startNode == endNode && startNode != -1 {
		fmt.Printf("Powerlines must be between nodes and cannot (for a single line) loop\n")
		return nil
	} else if startNode != -1 && endNode != -1 {
		// This might be a duplicate line.
		cost := p.grid.Cost(startNode, endNode)
		if cost != -1 {
			fmt.Printf("There already is a line from %v to %v.\n", startNode, endNode)
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
	p.powerLines[p.nextPowerLine] = &line
	p.nextPowerLine++

	return &line
}

func (p *PowerGrid) IteratePlants(iterate func(*PowerPlant)) {
	for _, plant := range p.powerPlants {
		iterate(plant)
	}
}

func (p *PowerGrid) IterateLines(iterate func(*PowerLine)) {
	for _, line := range p.powerLines {
		iterate(line)
	}
}
