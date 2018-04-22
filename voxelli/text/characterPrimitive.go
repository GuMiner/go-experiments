package text

import (
	"go-experiments/common/config"
	"go-experiments/voxelli/utils"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const quadVertLength = 6

var ccwQuadVert = []mgl32.Vec3{
	mgl32.Vec3{0, 1, 0},
	mgl32.Vec3{1, 0, 0},
	mgl32.Vec3{1, 1, 0},

	mgl32.Vec3{0, 0, 0},
	mgl32.Vec3{1, 0, 0},
	mgl32.Vec3{0, 1, 0}}

var cwQuadVert = []mgl32.Vec3{
	mgl32.Vec3{0, 1, 0},
	mgl32.Vec3{1, 1, 0},
	mgl32.Vec3{1, 0, 0},

	mgl32.Vec3{0, 0, 0},
	mgl32.Vec3{0, 1, 0},
	mgl32.Vec3{1, 0, 0}}

var ccwQuadUv = []mgl32.Vec2{
	mgl32.Vec2{0, 0},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{1, 0},

	mgl32.Vec2{0, 1},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{0, 0}}

var cwQuadUv = []mgl32.Vec2{
	mgl32.Vec2{1, 0},
	mgl32.Vec2{0, 0},
	mgl32.Vec2{0, 1},

	mgl32.Vec2{1, 1},
	mgl32.Vec2{1, 0},
	mgl32.Vec2{0, 1}}

func computeCharacterScale(textureScale utils.IntVec2) (float32, float32) {
	return float32(textureScale.X()) * commonConfig.Config.Text.PixelsToVerticesScale,
		float32(textureScale.Y()) * commonConfig.Config.Text.PixelsToVerticesScale
}

func generateCharacterPrimitive(
	offset float32,
	textureOffset, textureScale utils.IntVec2,
	textureSize int32,
	flip bool) ([]mgl32.Vec3, []mgl32.Vec2, float32) {

	xScale, yScale := computeCharacterScale(textureScale)

	positionBuffer := make([]mgl32.Vec3, quadVertLength)
	for i := 0; i < len(positionBuffer); i++ {
		var xVert, yVert float32
		if flip {
			xVert = cwQuadVert[i].X()
			yVert = cwQuadVert[i].Y()
		} else {
			xVert = ccwQuadVert[i].X()
			yVert = ccwQuadVert[i].Y()
		}

		positionBuffer[i] = mgl32.Vec3{xVert*xScale + offset, yVert * yScale, 0}
	}

	texturePositionBuffer := make([]mgl32.Vec2, len(ccwQuadUv))
	for i := 0; i < len(positionBuffer); i++ {
		x := textureOffset.X()

		var uvX, uvY float32
		if flip {
			uvX = cwQuadUv[i].X()
			uvY = cwQuadUv[i].Y()
		} else {
			uvX = ccwQuadUv[i].X()
			uvY = ccwQuadUv[i].Y()
		}

		if uvX > 0.5 {
			x += textureScale.X()
		}

		y := textureOffset.Y()
		if uvY > 0.5 {
			y += textureScale.Y()
		}

		texturePositionBuffer[i] = mgl32.Vec2{float32(x) / float32(textureSize), float32(y) / float32(textureSize)}
	}

	return positionBuffer, texturePositionBuffer, xScale
}

func renderPrimitive() {
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(ccwQuadVert)))
}
