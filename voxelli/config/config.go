package config

import (
	"encoding/json"
	"fmt"
	"go-experiments/common/config"
	"go-experiments/common/io"

	"github.com/go-gl/mathgl/mgl32"
)

type Agent struct {
	NeuralNet      []int
	WiggleDistance float32
	WiggleSpeed    float32
}

type Evolver struct {
	MaxGenerationLifetime float32
	ReportInterval        int
	SpeedCheckTime        float32
	MutationAmount        float32
	SelectionPercent      float32
	MutationProbability   float32
	CrossoverProbability  float32
	Mode                  string
}

type Vehicle struct {
	MinVelocity       float32
	MaxVelocity       float32
	AccelerationScale float32
	SteeringScale     float32
}

type Simulation struct {
	Agent          Agent
	Evolver        Evolver
	Vehicle        Vehicle
	PopulationSize int
	RoadwaySize    int
}

type Text struct {
	RuneFontSize          int
	BorderSize            int
	PixelsToVerticesScale float32
	FontFile              string
}

type ColorGradient struct {
	Steps      int
	Saturation float32
	Luminosity float32
}

type Shadows struct {
	Projection commonConfig.Projection
	Position   commonConfig.SerializableVec3
	Forwards   commonConfig.SerializableVec3
	Up         commonConfig.SerializableVec3
}

type DefaultCamera struct {
	Position commonConfig.SerializableVec3
	Forwards commonConfig.SerializableVec3
	Up       commonConfig.SerializableVec3
}

type Camera struct {
	MotionSpeed   float32
	RotationSpeed float32
	Default       DefaultCamera
}

type Configuration struct {
	Simulation    Simulation
	Window        commonConfig.Window
	Text          Text
	ColorGradient ColorGradient
	Shadows       Shadows
	Perspective   commonConfig.Perspective
	Camera        Camera
}

var Config Configuration

func (c *Camera) GetDefaultPos() mgl32.Vec3 {
	return mgl32.Vec3{
		c.Default.Position.X,
		c.Default.Position.Y,
		c.Default.Position.Z}
}

func (c *Camera) GetDefaultForwards() mgl32.Vec3 {
	return mgl32.Vec3{
		c.Default.Forwards.X,
		c.Default.Forwards.Y,
		c.Default.Forwards.Z}
}

func (c *Camera) GetDefaultUp() mgl32.Vec3 {
	return mgl32.Vec3{
		c.Default.Up.X,
		c.Default.Up.Y,
		c.Default.Up.Z}
}

func Load(configFileName string) {
	bytes := commonIo.ReadFileAsBytes(configFileName)

	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}

	fmt.Printf("Read in config '%v'.\n", configFileName)
	fmt.Printf("  Config data: %v\n\n", Config)
}
