package flat

import (
	"go-experiments/sim/engine/power"
	"go-experiments/sim/ui/lines"
	"go-experiments/sim/ui/region"

	"github.com/go-gl/mathgl/mgl32"
)

// TODO: Clearly these won't scale up, but they get a start with visuals so the focus can be on gameplay.
func RenderPowerPlants(grid *power.PowerGrid, camera *Camera, shadingProgram *region.RegionShaderProgram) {
	shadingProgram.PreRender()

	grid.IteratePlants(func(plant *power.PowerPlant) {
		region := plant.GetRegion()
		mappedRegion := camera.MapEngineRegionToScreen(region)
		shadingProgram.Render(mappedRegion, mgl32.Vec3{0.5, 0.5, 0.0})
	})
}

func RenderPowerLines(grid *power.PowerGrid, camera *Camera, shadingProgram *lines.LinesShaderProgram) {
	shadingProgram.PreRender()

	lines := make([][2]mgl32.Vec2, 0)
	grid.IterateLines(func(line *power.PowerLine) {
		mappedLine := camera.MapEngineLineToScreen(line.GetLine())
		lines = append(lines, mappedLine)
		if len(lines) > 200 {
			shadingProgram.Render(lines, mgl32.Vec3{0, 1, 0})
			lines = make([][2]mgl32.Vec2, 0)
		}
	})

	if len(lines) > 0 {
		shadingProgram.Render(lines, mgl32.Vec3{0, 1, 0})
	}
}
