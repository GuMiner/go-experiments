package subtile

import "github.com/go-gl/mathgl/mgl32"

// Defines common logic used heavily in subtile manipulation.
func GetRegionIndices(pos mgl32.Vec2, regionSize int) (x, y int) {
	x = int(pos.X()) / regionSize
	if pos.X() < 0 {
		x--
	}

	y = int(pos.Y()) / regionSize
	if pos.Y() < 0 {
		y--
	}

	return x, y
}

func GetLocalIndices(pos mgl32.Vec2, regionX, regionY, regionSize int) (x, y int) {
	x = int(pos.X()) - regionX*regionSize
	if pos.X() < 0 {
		x--
	}

	y = int(pos.Y()) - regionY*regionSize
	if pos.Y() < 0 {
		y--
	}

	return x, y
}

func GetLocalFloatIndices(pos mgl32.Vec2, regionX, regionY, regionSize int) mgl32.Vec2 {
	return pos.Sub(mgl32.Vec2{
		float32(regionX), float32(regionY)}.Mul(
		float32(regionSize)))
}
