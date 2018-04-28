package power

import "go-experiments/sim/config"

type PowerPlantType int
type PowerPlantSize int

// Ordering is important! Reordering these changes their numerical values
const (
	NaturalGas PowerPlantType = iota
	Coal
	Geothermal
	Nuclear
	Wind
	Solar
)

const (
	Small PowerPlantSize = iota
	Large
)

func GetPowerOutputAndSize(plantType PowerPlantType, plantSize PowerPlantSize) (output int, size int) {
	var plant config.PowerPlant

	switch plantType {
	case NaturalGas:
		plant = config.Config.Power.NaturalGas
	case Coal:
		plant = config.Config.Power.Coal
	case Geothermal:
		plant = config.Config.Power.Geothermal
	case Nuclear:
		plant = config.Config.Power.Nuclear
	case Wind:
		plant = config.Config.Power.Wind
	default: // Solar
		plant = config.Config.Power.Solar
	}

	switch plantSize {
	case Small:
		output = plant.SmallOutput
		size = plant.SmallSize
	default: // Large
		output = plant.LargeOutput
		size = plant.LargeSize
	}

	return output, size
}
