package commonColor

import (
	color "github.com/gerow/go-color"
	"github.com/go-gl/mathgl/mgl32"
)

var colorGradient []mgl32.Vec3

// Initializes the color gradient for future lookups
func InitializeColorGradient(maxIterations int, saturation float32, luminosity float32) {
	colorGradient = make([]mgl32.Vec3, maxIterations)

	for idx := range colorGradient {
		hue := float32(idx) / (float32(maxIterations) * 1.15) // We offset this a bit so we end solidly in purple and not back in red.
		color := color.HSL{H: float64(hue), S: float64(saturation), L: float64(luminosity)}.ToRGB()
		colorGradient[idx] = mgl32.Vec3{float32(color.R), float32(color.G), float32(color.B)}
	}
}

// Looks up a color in the color gradient by percent, 0 == red, 1 == purple
func LookupColor(percent float32) mgl32.Vec3 {
	id := int(percent * float32(len(colorGradient)))

	if id < 0 {
		id = 0
	} else if id >= len(colorGradient) {
		id = len(colorGradient) - 1
	}

	return colorGradient[id]
}
