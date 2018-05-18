package commonShadow

import (
	"fmt"
	"go-experiments/common/commonconfig"
	"go-experiments/common/commonmath"
	"go-experiments/common/commonopengl"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type ShadowBuffer struct {
	shadowBuffer  uint32
	shadowTexture uint32

	width  int32
	height int32
}

func NewShadowBuffer() *ShadowBuffer {
	var currentDrawBuffer int32
	gl.GetIntegerv(gl.DRAW_BUFFER, &currentDrawBuffer)

	var buffer ShadowBuffer

	maxTextureSize := commonMath.MinInt32(2048, commonOpenGl.GetGlCaps().MaxTextureSize)
	buffer.width = maxTextureSize
	buffer.height = maxTextureSize

	gl.GenTextures(1, &buffer.shadowTexture)
	gl.BindTexture(gl.TEXTURE_2D, buffer.shadowTexture)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.DEPTH_COMPONENT32, buffer.width, buffer.height)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_COMPARE_MODE, gl.COMPARE_REF_TO_TEXTURE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_COMPARE_FUNC, gl.LEQUAL)

	gl.GenFramebuffers(1, &buffer.shadowBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, buffer.shadowBuffer)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, buffer.shadowTexture, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)

	framebufferStatus := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if framebufferStatus != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("The shadow framebuffer status is not complete, found: %v\n", framebufferStatus))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return &buffer
}

func (r *ShadowBuffer) GetTextureId() uint32 {
	return r.shadowTexture
}

func (r *ShadowBuffer) Delete() {
	gl.DeleteTextures(1, &r.shadowTexture)
	gl.DeleteFramebuffers(1, &r.shadowBuffer)
}

func (r *ShadowBuffer) PrepareCamera() (mgl32.Mat4, mgl32.Mat4) {
	gl.Viewport(0, 0, r.width, r.height)

	projection := mgl32.Ortho(
		commonConfig.Config.Shadows.Projection.Left,
		commonConfig.Config.Shadows.Projection.Right,
		commonConfig.Config.Shadows.Projection.Bottom,
		commonConfig.Config.Shadows.Projection.Top,
		commonConfig.Config.Shadows.Projection.Near,
		commonConfig.Config.Shadows.Projection.Far)

	position := mgl32.Vec3{
		commonConfig.Config.Shadows.Position.X,
		commonConfig.Config.Shadows.Position.Y,
		commonConfig.Config.Shadows.Position.Z}

	cameraMatrix := mgl32.LookAtV(
		position,
		position.Add(
			mgl32.Vec3{
				commonConfig.Config.Shadows.Forwards.X,
				commonConfig.Config.Shadows.Forwards.Y,
				commonConfig.Config.Shadows.Forwards.Z}),
		mgl32.Vec3{
			commonConfig.Config.Shadows.Up.X,
			commonConfig.Config.Shadows.Up.Y,
			commonConfig.Config.Shadows.Up.Z})

	return projection, cameraMatrix
}

func (r *ShadowBuffer) RenderToBuffer(renderFunction func()) {
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.shadowBuffer)

	renderFunction()

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
