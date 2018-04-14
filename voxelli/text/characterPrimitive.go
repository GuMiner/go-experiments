package text

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var ccwQuadVert = []mgl32.Vec3{
	mgl32.Vec3{-50, 50, 40},
	mgl32.Vec3{50, -50, 40},
	mgl32.Vec3{50, 50, 40},

	mgl32.Vec3{-50, -50, 40},
	mgl32.Vec3{50, -50, 40},
	mgl32.Vec3{-50, 50, 40}}

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

func sendPrimitivesToDevice(positionVbo, colorVbo, texPosVbo uint32) {
	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BindBuffer(gl.ARRAY_BUFFER, positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadVert)*3*4, gl.Ptr(ccwQuadVert), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, colorVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadColor)*3*4, gl.Ptr(ccwQuadColor), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, texPosVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadUv)*2*4, gl.Ptr(ccwQuadUv), gl.STATIC_DRAW)
}
