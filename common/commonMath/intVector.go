package commonMath

import (
	"github.com/go-gl/mathgl/mgl32"
)

type IntVec2 [2]int
type IntVec3 [3]int
type Color [4]uint8

func (v IntVec2) X() int {
	return v[0]
}

func (v IntVec2) Y() int {
	return v[1]
}

func (v IntVec3) X() int {
	return v[0]
}

func (v IntVec3) Y() int {
	return v[1]
}

func (v IntVec3) Z() int {
	return v[2]
}

func (v IntVec3) AsFloatVector() mgl32.Vec3 {
	var pos mgl32.Vec3
	pos[0] = float32(v.X())
	pos[1] = float32(v.Y())
	pos[2] = float32(v.Z())
	return pos
}

func (v Color) R() uint8 {
	return v[0]
}

func (v Color) G() uint8 {
	return v[1]
}

func (v Color) B() uint8 {
	return v[2]
}

func (v Color) A() uint8 {
	return v[3]
}

func (v Color) AsFloatVector() mgl32.Vec4 {
	var colors mgl32.Vec4
	colors[0] = float32(v.R()) / 255.0
	colors[1] = float32(v.G()) / 255.0
	colors[2] = float32(v.B()) / 255.0
	colors[3] = float32(v.A()) / 255.0
	return colors
}

func (v Color) AsOpaqueFloatVector() mgl32.Vec3 {
	var colors mgl32.Vec3
	colors[0] = float32(v.R()) / 255.0
	colors[1] = float32(v.G()) / 255.0
	colors[2] = float32(v.B()) / 255.0
	return colors
}

func NewColor(bgr uint32) Color {
	r := uint8(bgr)
	g := uint8(bgr >> 8)
	b := uint8(bgr >> 16)
	a := uint8(bgr >> 24)
	return Color{r, g, b, a}
}
