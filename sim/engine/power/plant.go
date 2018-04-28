package power

import (
	"github.com/go-gl/mathgl/mgl32"
)

type PowerPlant struct {
	location    mgl32.Vec2
	size        float32 // All plants are assumed square
	orientation float32

	output int // kW
	// TODO: Add capacity factor, things impacted by, required resources, etc.
}
