package overlay

import (
	"go-experiments/common/commonopengl"
	"go-experiments/sim/config"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var ccwQuad = []mgl32.Vec2{
	mgl32.Vec2{-1, 1},
	mgl32.Vec2{1, -1},
	mgl32.Vec2{1, 1},

	mgl32.Vec2{-1, -1},
	mgl32.Vec2{1, -1},
	mgl32.Vec2{-1, 1}}

type OverlayShaderProgram struct {
	program uint32

	zOrderLoc       int32
	offsetLoc       int32
	scaleLoc        int32
	overlayImageLoc int32
	terrainSizeLoc  int32

	vao uint32
	vbo uint32
}

func NewOverlayShaderProgram() *OverlayShaderProgram {
	overlay := OverlayShaderProgram{}

	// Setup the OpenGL program
	overlay.program = commonOpenGl.CreateProgram("./ui/shaders/overlay")

	overlay.zOrderLoc = gl.GetUniformLocation(overlay.program, gl.Str("zOrder\x00"))
	overlay.offsetLoc = gl.GetUniformLocation(overlay.program, gl.Str("offset\x00"))
	overlay.scaleLoc = gl.GetUniformLocation(overlay.program, gl.Str("scale\x00"))
	overlay.overlayImageLoc = gl.GetUniformLocation(overlay.program, gl.Str("overlayImage\x00"))
	overlay.terrainSizeLoc = gl.GetUniformLocation(overlay.program, gl.Str("terrainSize\x00"))

	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &overlay.vao)
	gl.BindVertexArray(overlay.vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &overlay.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, overlay.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	// 2 -- 2 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuad)*2*4, gl.Ptr(ccwQuad), gl.STATIC_DRAW)

	return &overlay
}

func (shaderProgram *OverlayShaderProgram) PreRender() {
	gl.UseProgram(shaderProgram.program)
	gl.BindVertexArray(shaderProgram.vao)
	gl.Uniform1f(shaderProgram.terrainSizeLoc, float32(config.Config.Terrain.RegionSize))
}

func (shaderProgram *OverlayShaderProgram) Render(overlay *Overlay) {
	// Setup texture
	gl.ActiveTexture(gl.TEXTURE0 + overlay.textureBindLocation)
	gl.BindTexture(gl.TEXTURE_2D, overlay.textureId)
	gl.Uniform1i(shaderProgram.overlayImageLoc, int32(overlay.textureBindLocation)-gl.TEXTURE0)

	// Setup location
	gl.Uniform1f(shaderProgram.zOrderLoc, overlay.zOrder)
	gl.Uniform2f(shaderProgram.offsetLoc, overlay.offset.X(), overlay.offset.Y())
	gl.Uniform2f(shaderProgram.scaleLoc, overlay.scale.X(), overlay.scale.Y())

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(ccwQuad)))
}

func (shaderProgram *OverlayShaderProgram) Delete() {
	gl.DeleteBuffers(1, &shaderProgram.vbo)
	gl.DeleteVertexArrays(1, &shaderProgram.vao)
	gl.DeleteProgram(shaderProgram.program)
}
