package config

import (
	"encoding/json"
	"fmt"
	"go-experiments/voxelli/utils"

	"github.com/go-gl/mathgl/mgl32"
)

type SerializableVec3 struct {
	X float32
	Y float32
	Z float32
}

type Projection struct {
	Left   float32
	Right  float32
	Bottom float32
	Top    float32
	Near   float32
	Far    float32
}

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

type Window struct {
	Width           int
	Height          int
	Title           string
	BackgroundColor SerializableVec3
	OpenGlMajor     int
	OpenGlMinor     int
}

type ColorGradient struct {
	Steps      int
	Saturation float32
	Luminosity float32
}

type Shadows struct {
	Projection Projection
	Position   SerializableVec3
	Forwards   SerializableVec3
	Up         SerializableVec3
}

type Perspective struct {
	FovY float32
	Near float32
	Far  float32
}

type DefaultCamera struct {
	Position SerializableVec3
	Forwards SerializableVec3
	Up       SerializableVec3
}

type Camera struct {
	MotionSpeed   float32
	RotationSpeed float32
	Default       DefaultCamera
}

type Configuration struct {
	Simulation    Simulation
	Window        Window
	Text          Text
	ColorGradient ColorGradient
	Shadows       Shadows
	Perspective   Perspective
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
	bytes := utils.ReadFileAsBytes(configFileName)

	if err := json.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}

	fmt.Printf("Read in config '%v'.\n", configFileName)
	fmt.Printf("  Config data: %v\n\n", Config)
}
