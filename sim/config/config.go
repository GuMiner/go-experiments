package config

import (
	"encoding/json"
	"fmt"
	"go-experiments/common/commonconfig"
	"go-experiments/common/commonio"
)

type PowerPlant struct {
	SmallOutput int
	SmallSize   int

	LargeOutput int
	LargeSize   int
}

type Power struct {
	Coal       PowerPlant
	Nuclear    PowerPlant
	NaturalGas PowerPlant
	Wind       PowerPlant
	Solar      PowerPlant
	Geothermal PowerPlant
}

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

type TerrainUi struct {
	// Colors range from [0-256), not [0 - 1)
	WaterColor commonConfig.SerializableVec3
	SandColor  commonConfig.SerializableVec3
	GrassColor commonConfig.SerializableVec3
	HillColor  commonConfig.SerializableVec3
	RockColor  commonConfig.SerializableVec3
	SnowColor  commonConfig.SerializableVec3
}

type CameraConfig struct {
	MouseScrollFactor float32
	KeyMotionFactor   float32
}

type Ui struct {
	TerrainUi TerrainUi
	Camera    CameraConfig
}

type DrawConfig struct {
	SnapNodeCount       int
	MinSnapNodeDistance float32
}

type SnapConfig struct {
	SnapAngleDivision  int
	SnapGridResolution int
}

type Configuration struct {
	Terrain Terrain
	Power   Power
	Ui      Ui
	Draw    DrawConfig
	Snap    SnapConfig
}

var Config Configuration

func Load(configFileName string, commonConfigFileName string) {
	commonConfig.Load(commonConfigFileName)
	bytes := commonIo.ReadFileAsBytes(configFileName)

	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}

	fmt.Printf("Read in config '%v'.\n", configFileName)
	fmt.Printf("  Config data: %v\n\n", Config)
}
