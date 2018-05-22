package flat

import (
	"go-experiments/common/commonmath"
	"go-experiments/sim/engine"
	"go-experiments/sim/ui/region"

	"github.com/go-gl/mathgl/mgl32"
)

func RenderSnapNodes(snapElements *engine.SnapElements, camera *Camera, shadingProgram *region.RegionShaderProgram) {
	shadingProgram.PreRender()

	snappedNode := snapElements.GetSnappedNode()
	if snappedNode != nil {
		region := commonMath.Region{
			RegionType:  commonMath.CircleRegion,
			Position:    snappedNode.Element.GetSnapNodes()[snappedNode.SnapNodeIdx],
			Scale:       50,
			Orientation: 0}

		mappedRegion := camera.MapEngineRegionToScreen(&region)
		color := mgl32.Vec3{0.0, 1.0, 0.0}
		shadingProgram.Render(mappedRegion, color)
	}
}
