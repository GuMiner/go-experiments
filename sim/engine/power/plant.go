package power

import (
	"fmt"

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

type PowerPlants struct {
	powerPlants    map[int64]PowerPlant
	nextPowerPlant int64
}

func NewPowerPlants() *PowerPlants {
	return &PowerPlants{
		powerPlants:    make(map[int64]PowerPlant),
		nextPowerPlant: 0}
}

func (p *PowerPlants) Add(pos mgl32.Vec2, plantType PowerPlantType, plantSize PowerPlantSize) {
	output, size := GetPowerOutputAndSize(plantType, plantSize)

	p.powerPlants[p.nextPowerPlant] = PowerPlant{
		id:          p.nextPowerPlant,
		location:    pos,
		plantType:   plantType,
		namedSize:   plantSize,
		size:        float32(size),
		orientation: 0, // TODO: Rotation
		output:      output}
	p.nextPowerPlant++

	fmt.Printf("Added power plant '%v'.\n", p.powerPlants[p.nextPowerPlant-1])
}
