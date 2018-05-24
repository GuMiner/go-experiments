package config

type GenerationParameters struct {
	Seed int

	// The inverse of the noise scale. Lower values == more granular.
	MaxNoiseScale float32
	MedNoiseScale float32
	MinNoiseScale float32

	// Amount each noise amount contributes to the final total
	// Can add up to anything.
	MaxNoiseContribution float32
	MedNoiseContribution float32
	MinNoiseContribution float32

	// How much the lower end of the noise spectrum is flattened
	PowerFactor float32
}

type Terrain struct {
	// Levels for which the given terrain begins.
	WaterLevel float32
	SandLevel  float32
	GrassLevel float32
	HillLevel  float32
	RockLevel  float32
	SnowLevel  float32

	Generation GenerationParameters
	RegionSize int
}
