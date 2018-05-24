package config

type PowerPlant struct {
	SmallOutput int
	SmallSize   int

	LargeOutput int
	LargeSize   int

	Cost float32
}

type Power struct {
	PowerPlantTypes map[string]PowerPlant
	PowerLineCost   float32 // Cost per unit

	// Generated at run-time as ordering of maps is not guaranteed
	IdToNameMap map[int]string
}
