package flat

import (
	"go-experiments/sim/engine/road"
	"go-experiments/sim/ui/lines"

	"github.com/go-gl/mathgl/mgl32"
)

func RenderRoadLines(grid *road.RoadGrid, camera *Camera, shadingProgram *lines.LinesShaderProgram) {
	shadingProgram.PreRender()

	lines := make([][2]mgl32.Vec2, 0)
	grid.IterateLines(func(line *road.RoadLine) {
		mappedLine := camera.MapEngineLineToScreen(line.GetLine())
		lines = append(lines, mappedLine)
		if len(lines) > 200 {
			shadingProgram.Render(lines, mgl32.Vec3{1, 0, 0})
			lines = make([][2]mgl32.Vec2, 0)
		}
	})

	if len(lines) > 0 {
		shadingProgram.Render(lines, mgl32.Vec3{1, 0, 0})
	}
}
