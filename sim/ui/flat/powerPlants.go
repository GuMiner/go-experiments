package flat

import (
	"go-experiments/sim/engine/power"
	"go-experiments/sim/ui/region"

	"github.com/go-gl/mathgl/mgl32"
)

func RenderPowerPlants(plants *power.PowerPlants, camera *Camera, shadingProgram *region.RegionShaderProgram) {
	shadingProgram.PreRender()

	plants.Iterate(func(plant *power.PowerPlant) {
		region := plant.GetRegion()
		mappedRegion := camera.MapEngineRegionToScreen(region)
		shadingProgram.Render(&mappedRegion, mgl32.Vec3{0.5, 0.5, 0.0})
	})
}
