package power

import (
	"go-experiments/sim/config"
	"go-experiments/sim/input/editorEngine"
)

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

func GetPlantType(plantAddMode editorEngine.EditorPlantAddMode) PowerPlantType {
	var plantType PowerPlantType

	switch plantAddMode {
	case editorEngine.CoalPlant:
		plantType = Coal
	case editorEngine.GeothermalPlant:
		plantType = Geothermal
	case editorEngine.NaturalGasPlant:
		plantType = NaturalGas
	case editorEngine.NuclearPlant:
		plantType = Nuclear
	case editorEngine.WindPlant:
		plantType = Wind
	default: // Solar
		plantType = Solar
	}

	return plantType
}

func GetPlantCost(plantType PowerPlantType) float32 {
	var cost float32
	switch plantType {
	case Coal:
		cost = config.Config.Power.Coal.Cost
	case Geothermal:
		cost = config.Config.Power.Geothermal.Cost
	case NaturalGas:
		cost = config.Config.Power.NaturalGas.Cost
	case Nuclear:
		cost = config.Config.Power.Nuclear.Cost
	case Wind:
		cost = config.Config.Power.Wind.Cost
	default: // Solar
		cost = config.Config.Power.Solar.Cost
	}

	return cost
}
