package config

import "go-experiments/sim/engine/resource"

type Building struct {
	Name string
	Size int
	Cost float32

	// Required resource for each day for the building to function
	RequiredBasics resource.UtilityResource

	// Required inputs for the building to perform its function and output product.
	Inputs  []resource.ResourceAmount
	Outputs []resource.ResourceAmount

	// How much of each resource the building can store. Includes inputs and outputs.
	StorageCapacity map[string]float32
}
