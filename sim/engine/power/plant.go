package power

import (
	"github.com/go-gl/mathgl/mgl32"
)

type PowerPlant struct {
	id          int64
	location    mgl32.Vec2
	size        float32 // All plants are assumed square
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
