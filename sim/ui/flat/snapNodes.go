package flat

import (
	"go-experiments/common/commonmath"
	"go-experiments/sim/config"
	"go-experiments/sim/engine/element"
	"go-experiments/sim/ui/region"

	"github.com/go-gl/mathgl/mgl32"
)

func RenderSnapNodes(boardPos mgl32.Vec2, elementFinder *element.ElementFinder, camera *Camera, shadingProgram *region.RegionShaderProgram) {
	shadingProgram.PreRender()

	elements := elementFinder.KNearest(boardPos, config.Config.Draw.SnapNodeCount)
	for i, elem := range elements {
		if elem.Distance < config.Config.Draw.MinSnapNodeDistance {
			for j, snapNode := range elem.Element.GetSnapNodes() {
				region := commonMath.Region{
					RegionType:  commonMath.CircleRegion,
					Position:    snapNode,
					Scale:       50,
					Orientation: 0}
				mappedRegion := camera.MapEngineRegionToScreen(region)
				color := mgl32.Vec3{1.0, 0.0, 1.0}
				if i == 0 && j == 0 {
					color = mgl32.Vec3{0.0, 1.0, 0.0}
				}
				shadingProgram.Render(&mappedRegion, color)
			}
		}
	}
}
