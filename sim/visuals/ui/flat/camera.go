package flat

import (
	"go-experiments/common/opengl"

	"go-experiments/sim/config"
	"go-experiments/sim/input"
	"go-experiments/voxelli/utils"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	zoomFactor float32
	offset     mgl32.Vec2
}

func NewCamera() *Camera {
	camera := Camera{
		zoomFactor: 0.0,
		offset:     mgl32.Vec2{0, 0}}

	return &camera
}

func (c *Camera) Update(frameTime float32) {
	if input.ScrollEvent {
		c.zoomFactor += input.MouseScrollOffset[0] * config.Config.Ui.Camera.MouseScrollFactor
		input.ScrollEvent = false
	}

	keyMotionAmount := frameTime * config.Config.Ui.Camera.KeyMotionFactor * c.getScaleMotionFactor()
	if input.IsPressed(input.MoveLeft) {
		c.offset[0] -= keyMotionAmount
	}

	if input.IsPressed(input.MoveRight) {
		c.offset[0] += keyMotionAmount
	}

	if input.IsPressed(input.MoveUp) {
		c.offset[1] -= keyMotionAmount
	}

	if input.IsPressed(input.MoveDown) {
		c.offset[1] += keyMotionAmount
	}
}

func (c *Camera) ComputeVisibleRegions() []utils.IntVec2 {
	windowSize := commonOpenGl.GetWindowSize()
	regionSize := config.Config.Terrain.RegionSize

	// TODO: Handle zoom factor
	minTile := c.MapToBoard(mgl32.Vec2{0, 0}).Mul(1.0 / float32(regionSize))
	maxTile := c.MapToBoard(windowSize).Mul(1.0 / float32(regionSize))

	visibleTiles := make([]utils.IntVec2, 0)
	for i := int(minTile.X() - 1.0); i <= int(maxTile.X()+1.0); i++ {
		for j := int(minTile.Y() - 1.0); j <= int(maxTile.Y()+1.0); j++ {
			visibleTiles = append(visibleTiles, utils.IntVec2{i, j})
		}
	}

	return visibleTiles
}

func (c *Camera) MapToBoard(screenPos mgl32.Vec2) mgl32.Vec2 {
	return screenPos.Add(c.offset)
}

// Resizes a full-size region to the appropriate scale given the current screen size and zoom factor
func (c *Camera) GetRegionScale() mgl32.Vec2 {
	// TODO: Handle zoom factor

	regionSize := config.Config.Terrain.RegionSize
	windowSize := commonOpenGl.GetWindowSize()
	return mgl32.Vec2{float32(regionSize) / windowSize.X(), float32(regionSize) / windowSize.Y()}
}

func (c *Camera) GetRegionOffset(x, y int) mgl32.Vec2 {
	regionSize := config.Config.Terrain.RegionSize
	windowSize := commonOpenGl.GetWindowSize()

	regionStart := mgl32.Vec2{float32(x * regionSize), float32(y * regionSize)}
	regionPos := regionStart.Sub(c.offset)
	return mgl32.Vec2{regionPos.X() / windowSize.X(), regionPos.Y() / windowSize.Y()}
}

func (c *Camera) getScaleMotionFactor() float32 {
	// TODO:
	return 1.0 // c.zoomFactor
}
