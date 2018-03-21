package main

type IntVec3 [3]int
type Color [3]uint8

func (v IntVec3) X() int {
	return v[0]
}

func (v IntVec3) Y() int {
	return v[1]
}

func (v IntVec3) Z() int {
	return v[2]
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

func NewColor(bgr uint32) Color {
	r := uint8(bgr)
	g := uint8(bgr >> 8)
	b := uint8(bgr >> 16)
	return Color{r, g, b}
}
