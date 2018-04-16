package text

import (
	"go-experiments/voxelli/utils"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var ccwQuadVert = []mgl32.Vec3{
	mgl32.Vec3{0, 1, 0},
	mgl32.Vec3{1, 0, 0},
	mgl32.Vec3{1, 1, 0},

	mgl32.Vec3{0, 0, 0},
	mgl32.Vec3{1, 0, 0},
	mgl32.Vec3{0, 1, 0}}

var ccwQuadUv = []mgl32.Vec2{
	mgl32.Vec2{0, 0},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{1, 0},

	mgl32.Vec2{0, 1},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{0, 0}}

const pixelsToVerticesScale = 0.05 // Scales down the pixel size of a character to vertices

func sendPrimitivesToDevice(
	positionVbo, texPosVbo uint32,
	lastRuneOffset float32,
	textureOffset, textureScale utils.IntVec2,
	textureSize int32) float32 {

	xScale := float32(textureScale.X()) * pixelsToVerticesScale
	yScale := float32(textureScale.Y()) * pixelsToVerticesScale

	positionBuffer := make([]mgl32.Vec3, len(ccwQuadVert))
	for i := 0; i < len(positionBuffer); i++ {
		positionBuffer[i] = mgl32.Vec3{ccwQuadVert[i].X()*xScale + lastRuneOffset, ccwQuadVert[i].Y() * yScale, 0}
	}

	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BindBuffer(gl.ARRAY_BUFFER, positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(positionBuffer)*3*4, gl.Ptr(positionBuffer), gl.STATIC_DRAW)

	texturePositionBuffer := make([]mgl32.Vec2, len(ccwQuadUv))
	for i := 0; i < len(positionBuffer); i++ {
		x := textureOffset.X()
		if ccwQuadUv[i].X() > 0.5 {
			x += textureScale.X()
		}

		y := textureOffset.Y()
		if ccwQuadUv[i].Y() > 0.5 {
			y += textureScale.Y()
		}

		texturePositionBuffer[i] = mgl32.Vec2{float32(x) / float32(textureSize), float32(y) / float32(textureSize)}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, texPosVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(texturePositionBuffer)*2*4, gl.Ptr(texturePositionBuffer), gl.STATIC_DRAW)

	return xScale
}

func renderPrimitive() {
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(ccwQuadVert)))
}
