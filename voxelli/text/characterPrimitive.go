package text

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var ccwQuadVert = []mgl32.Vec3{
	mgl32.Vec3{-0.5, 0.5, 0},
	mgl32.Vec3{0.5, -0.5, 0},
	mgl32.Vec3{0.5, 0.5, 0},

	mgl32.Vec3{-0.5, -0.5, 0},
	mgl32.Vec3{0.5, -0.5, 0},
	mgl32.Vec3{-0.5, 0.5, 0}}

var ccwQuadColor = []mgl32.Vec3{
	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5},

	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5}}

var ccwQuadUv = []mgl32.Vec2{
	mgl32.Vec2{0, 0},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{1, 0},

	mgl32.Vec2{0, 1},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{0, 0}}

func sendPrimitivesToDevice(
	positionVbo, colorVbo, texPosVbo uint32,
	characterOffset mgl32.Vec2,
	minBounds, maxBounds mgl32.Vec2) {

	characterOffset3d := mgl32.Vec3{characterOffset.X(), characterOffset.Y(), 0}
	positionBuffer := make([]mgl32.Vec3, len(ccwQuadVert))
	for i := 0; i < len(positionBuffer); i++ {
		positionBuffer[i] = ccwQuadVert[i].Add(characterOffset3d)
	}

	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BindBuffer(gl.ARRAY_BUFFER, positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(positionBuffer)*3*4, gl.Ptr(positionBuffer), gl.STATIC_DRAW)

	// TODO: Regen color buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadColor)*3*4, gl.Ptr(ccwQuadColor), gl.STATIC_DRAW)

	texturePositionBuffer := make([]mgl32.Vec2, len(ccwQuadUv))
	for i := 0; i < len(positionBuffer); i++ {
		x := minBounds.X()
		if ccwQuadUv[i].X() > 0.5 {
			x = maxBounds.X()
		}

		y := minBounds.Y()
		if ccwQuadUv[i].Y() > 0.5 {
			y = maxBounds.Y()
		}

		texturePositionBuffer[i] = mgl32.Vec2{x, y}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, texPosVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(texturePositionBuffer)*2*4, gl.Ptr(texturePositionBuffer), gl.STATIC_DRAW)
}

func renderPrimitive() {
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(ccwQuadVert)))
}
