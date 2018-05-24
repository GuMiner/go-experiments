package config

import "github.com/go-gl/mathgl/mgl32"

type Resource struct {
	Name  string
	Units string
	Color mgl32.Vec2

	// What types of vehicles can carry this resource
	AllowedVehicleTypes []string

	// A scale factor for unit amounts.
	// For example, a vehicle of capacity 10 can carry 10*ResourceFactor units of this resource.
	ResourceFactor float32
}
