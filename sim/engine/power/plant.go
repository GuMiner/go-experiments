package power

import (
	"fmt"
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

type PowerPlant struct {
	id       int64
	location mgl32.Vec2

	plantType PowerPlantType
	namedSize PowerPlantSize

	size        float32 // All plants are assumed square, for now. (TODO)
	orientation float32

	output int // kW
	// TODO: Add capacity factor, things impacted by, required resources, etc.
}

// Implement Element
// Gets the central position of the element
func (p *PowerPlant) GetRegion() commonMath.Region {
	return commonMath.Region{
		RegionType:  commonMath.SquareRegion,
		Position:    p.location,
		Scale:       p.size,
		Orientation: p.orientation}
}

// Gets positions on the map that can be used to points of the element.
func (p *PowerPlant) GetSnapNodes() []mgl32.Vec2 {
	return make([]mgl32.Vec2, 0)
}

// Gets lines on the map that can be used to snap to *edges* of the element
func (p *PowerPlant) GetSnapEdges() [][2]mgl32.Vec2 {
	return make([][2]mgl32.Vec2, 0)
}

type PowerPlants struct {
	powerPlants    map[int64]*PowerPlant
	nextPowerPlant int64
}

func NewPowerPlants() *PowerPlants {
	return &PowerPlants{
		powerPlants:    make(map[int64]*PowerPlant),
		nextPowerPlant: 0}
}

func (p *PowerPlants) Add(pos mgl32.Vec2, plantType PowerPlantType, plantSize PowerPlantSize) *PowerPlant {
	output, size := GetPowerOutputAndSize(plantType, plantSize)

	plant := PowerPlant{
		id:          p.nextPowerPlant,
		location:    pos,
		plantType:   plantType,
		namedSize:   plantSize,
		size:        float32(size),
		orientation: 0, // TODO: Rotation
		output:      output}

	p.powerPlants[p.nextPowerPlant] = &plant
	fmt.Printf("Added power plant '%v'.\n", p.powerPlants[p.nextPowerPlant])

	p.nextPowerPlant++
	return &plant
}

func (p *PowerPlants) Iterate(iterate func(*PowerPlant)) {
	for _, plant := range p.powerPlants {
		iterate(plant)
	}
}
