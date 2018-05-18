package flat

import (
	"go-experiments/sim/engine/power"
	"go-experiments/sim/ui/region"

	"github.com/go-gl/mathgl/mgl32"
)

func RenderPowerPlants(grid *power.PowerGrid, camera *Camera, shadingProgram *region.RegionShaderProgram) {
	shadingProgram.PreRender()

	grid.IteratePlants(func(plant *power.PowerPlant) {
		region := plant.GetRegion()
		mappedRegion := camera.MapEngineRegionToScreen(region)
		shadingProgram.Render(&mappedRegion, mgl32.Vec3{0.5, 0.5, 0.0})
	})
}
