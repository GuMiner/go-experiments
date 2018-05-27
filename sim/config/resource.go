package config

import "github.com/go-gl/mathgl/mgl32"

type SpawnConfig struct {
	HeightRange mgl32.Vec2 // End of the map

	MinSpawnSize     float32
	SpawnProbability float32 // Set to 1 for this to spawn always in this height range.

	AmountPerArea             float32
	RegenerationPercentPerDay float32 // Set to 0 for finite resources
}

type Transformation struct {
	Resource    string  // The resource this resource transforms into
	Factor      float32 // A conversion factor, for unit conversions per day.
	DelayFactor int     // The number of days before this resource starts transforming.
}

type Resource struct {
	Name  string
	Color mgl32.Vec3

	Units string
	// A scale factor for unit amounts.
	// For example, a vehicle of capacity 10 can carry 10*ResourceFactor units of this resource.
	ResourceFactor float32

	// What types of vehicles can carry this resource
	AllowedVehicleTypes []string

	// How this resource spawns.
	SpawnSetup []SpawnConfig

	// How this resource mutates over time.
	Transformations []Transformation
}
