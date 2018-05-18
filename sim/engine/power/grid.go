package power

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type PowerGrid struct {
	grid           *ResizeableGraph
	powerPlants    map[int64]*PowerPlant
	nextPowerPlant int64
}

func NewPowerGrid() *PowerGrid {
	grid := PowerGrid{
		grid:           NewResizeableGraph(),
		powerPlants:    make(map[int64]*PowerPlant),
		nextPowerPlant: 0}

	return &grid
}

func (p *PowerGrid) Add(pos mgl32.Vec2, plantType PowerPlantType, plantSize PowerPlantSize) *PowerPlant {
	output, size := GetPowerOutputAndSize(plantType, plantSize)

	plant := PowerPlant{
		id:          p.nextPowerPlant,
		location:    pos,
		plantType:   plantType,
		namedSize:   plantSize,
		size:        float32(size),
		orientation: 0, // TODO: Rotation
		output:      output,
		gridId:      p.grid.AddNode()}

	p.powerPlants[p.nextPowerPlant] = &plant
	fmt.Printf("Added power plant '%v'.\n", p.powerPlants[p.nextPowerPlant])

	p.nextPowerPlant++

	return &plant
}

func (p *PowerGrid) IteratePlants(iterate func(*PowerPlant)) {
	for _, plant := range p.powerPlants {
		iterate(plant)
	}
}
