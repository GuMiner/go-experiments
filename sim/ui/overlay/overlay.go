package overlay

// Defines a full-screen quad
import (
	"github.com/go-gl/mathgl/mgl32"
)

type Overlay struct {
	offset mgl32.Vec2
	scale  mgl32.Vec2
	zOrder float32

	textureId           uint32
	textureBindLocation uint32
}

func NewOverlay() *Overlay {
	overlay := Overlay{
		offset:              mgl32.Vec2{0, 0},
		scale:               mgl32.Vec2{1, 1},
		zOrder:              float32(1.0),
		textureBindLocation: 0}

	return &overlay
}

func (overlay *Overlay) UpdateTexture(textureId uint32) {
	overlay.textureId = textureId
}

func (overlay *Overlay) UpdateLocation(offset mgl32.Vec2, scale mgl32.Vec2, zOrder float32) {
	overlay.offset = offset
	overlay.scale = scale
	overlay.zOrder = zOrder
}
