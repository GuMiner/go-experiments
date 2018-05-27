package config

import (
	"encoding/json"
	"fmt"
	"go-experiments/common/commonconfig"
	"go-experiments/common/commonio"
)

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

type SimConfig struct {
	SecondsPerDay   float32
	StartingSavings float32
	MaxDebt         float32
}

type Configuration struct {
	Terrain Terrain
	Power   Power

	Ui   Ui
	Draw DrawConfig
	Snap SnapConfig
	Sim  SimConfig

	Buildings []Building
	Resources []Resource
	Vehicles  []Vehicle
}

var Config Configuration

func loadSubConfigs(configFolder string) {
	// Sub-configs
	bytes := commonIo.ReadFileAsBytes(configFolder + "terrain.json")
	if err := json.Unmarshal(bytes, &Config.Terrain); err != nil {
		panic(err)
	}

	bytes = commonIo.ReadFileAsBytes(configFolder + "power.json")
	if err := json.Unmarshal(bytes, &Config.Power); err != nil {
		panic(err)
	}

	i := 0
	Config.Power.IdToNameMap = make(map[int]string)
	for key, _ := range Config.Power.PowerPlantTypes {
		Config.Power.IdToNameMap[i] = key
		i++
	}

	bytes = commonIo.ReadFileAsBytes(configFolder + "building.json")
	if err := json.Unmarshal(bytes, &Config.Buildings); err != nil {
		panic(err)
	}

	bytes = commonIo.ReadFileAsBytes(configFolder + "resource.json")
	if err := json.Unmarshal(bytes, &Config.Resources); err != nil {
		panic(err)
	}

	bytes = commonIo.ReadFileAsBytes(configFolder + "vehicle.json")
	if err := json.Unmarshal(bytes, &Config.Vehicles); err != nil {
		panic(err)
	}
}

func Load(configFolder string, commonConfigFileName string) {
	commonConfig.Load(commonConfigFileName)

	Config = Configuration{}

	bytes := commonIo.ReadFileAsBytes(configFolder + "config.json")
	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}

	loadSubConfigs(configFolder)

	fmt.Printf("Read in config from '%v'.\n", configFolder)
	fmt.Printf("  Config data: %v\n\n", Config)
}
