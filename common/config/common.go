package commonConfig

import (
	"encoding/json"
	"fmt"
	"go-experiments/common/io"
)

type SerializableVec3 struct {
	X float32
	Y float32
	Z float32
}

type CommonConfig struct {
	// Camera
	Perspective Perspective
	Window      Window

	// Color
	ColorGradient ColorGradient

	// Shadows
	Shadows Shadows

	// text
	Text Text
}

var Config CommonConfig

func Load(configFileName string) {
	bytes := commonIo.ReadFileAsBytes(configFileName)

	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}

	fmt.Printf("Read in common config '%v'.\n", configFileName)
	fmt.Printf("  Config data: %v\n\n", Config)
}
