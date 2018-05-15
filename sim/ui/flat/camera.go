package flat

import (
	"go-experiments/common/math"
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
		zoomFactor: 1.0,
		offset:     mgl32.Vec2{0, 0}}

	return &camera
}

func (c *Camera) Update(frameTime float32) {
	if input.ScrollEvent {
		scrollAmount := input.GetScrollOffset().Y()
		c.zoomFactor *= (1.0 + scrollAmount*config.Config.Ui.Camera.MouseScrollFactor)
		input.ScrollEvent = false
	}

	keyMotionAmount := frameTime * config.Config.Ui.Camera.KeyMotionFactor * (1.0 / c.zoomFactor)
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

func (c *Camera) getMinMaxVisibleRange() (minTile mgl32.Vec2, maxTile mgl32.Vec2) {
	regionSize := config.Config.Terrain.RegionSize

	minTile = c.MapToBoard(mgl32.Vec2{0, 0}).Mul(1.0 / float32(regionSize))
	maxTile = c.MapToBoard(mgl32.Vec2{1, 1}).Mul(1.0 / float32(regionSize))
	return minTile, maxTile
}

func (c *Camera) ComputeVisibleRegions() []utils.IntVec2 {
	minTile, maxTile := c.getMinMaxVisibleRange()

	visibleTiles := make([]utils.IntVec2, 0)
	for i := int(minTile.X() - 1.0); i <= int(maxTile.X()+1.0); i++ {
		for j := int(minTile.Y() - 1.0); j <= int(maxTile.Y()+1.0); j++ {
			visibleTiles = append(visibleTiles, utils.IntVec2{i, j})
		}
	}

	return visibleTiles
}

func (c *Camera) ComputePrecacheRegions() []utils.IntVec2 {
	minTile, maxTile := c.getMinMaxVisibleRange()

	visibleTiles := make([]utils.IntVec2, 0)
	for i := int(minTile.X() - 2.0); i <= int(maxTile.X()+2.0); i++ {
		for j := int(minTile.Y() - 2.0); j <= int(maxTile.Y()+2.0); j++ {
			if i == int(minTile.X()-2.0) ||
				i == int(minTile.X()+2.0) ||
				j == int(minTile.Y()-2.0) ||
				j == int(maxTile.Y()+2.0) {
				visibleTiles = append(visibleTiles, utils.IntVec2{i, j})
			}
		}
	}

	return visibleTiles
}

// Maps a position in pixels to the board
func (c *Camera) MapPixelPosToBoard(pixelPos mgl32.Vec2) mgl32.Vec2 {
	windowSize := commonOpenGl.GetWindowSize()
	return c.MapToBoard(mgl32.Vec2{pixelPos.X() / windowSize.X(), pixelPos.Y() / windowSize.Y()})
}

// Maps a (0, 0) to (1, 1) screen position to a board location.
func (c *Camera) MapToBoard(screenPos mgl32.Vec2) mgl32.Vec2 {
	windowSize := commonOpenGl.GetWindowSize()

	modifiedRegionPos := mgl32.Vec2{(screenPos.X() - 0.5) * windowSize.X(), (screenPos.Y() - 0.5) * windowSize.Y()}
	regionPos := modifiedRegionPos.Mul(1.0 / c.zoomFactor).Add(c.offset)

	return regionPos
}

// Maps a region on the board to a GLSL (-1, -1) to (1, 1) region
func (c *Camera) MapEngineRegionToScreen(region commonMath.Region) commonMath.Region {
	// The only variables that are updated (for now) are position and scale
	windowSize := commonOpenGl.GetWindowSize()

	region.Scale *= (1.0 / c.zoomFactor)
	region.Position = region.Position.Sub(c.offset).Mul(c.zoomFactor)
	region.Position = mgl32.Vec2{2 * region.Position.X() / windowSize.X(), -2 * region.Position.Y() / windowSize.Y()}
	return region
}

// Resizes a full-size region to the appropriate scale given the current screen size and zoom factor
// Returns the screen size (a full size tile will span from (0, 0) to (1, 1))
func (c *Camera) GetRegionScale() mgl32.Vec2 {
	regionSize := config.Config.Terrain.RegionSize
	windowSize := commonOpenGl.GetWindowSize()
	return mgl32.Vec2{
		c.zoomFactor * float32(regionSize) / windowSize.X(),
		c.zoomFactor * float32(regionSize) / windowSize.Y()}
}

// Returns the screen position ((0, 0) to (1, 1)) of the region tile requested
func (c *Camera) GetRegionOffset(x, y int) mgl32.Vec2 {
	regionSize := config.Config.Terrain.RegionSize
	windowSize := commonOpenGl.GetWindowSize()

	regionStart := mgl32.Vec2{float32(x * regionSize), float32(y * regionSize)}
	modifiedRegionStart := regionStart.Sub(c.offset).Mul(c.zoomFactor)

	return mgl32.Vec2{modifiedRegionStart.X()/windowSize.X() + 0.5, modifiedRegionStart.Y()/windowSize.Y() + 0.5}
}

func (c *Camera) getScaleMotionFactor() float32 {
	return c.zoomFactor
}
