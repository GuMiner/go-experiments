package main

// Defines a 1D rainbow texture
// Defines a full-screen quad
import (
	"github.com/gerow/go-color"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type RainbowTexture struct {
	openGlTextureSlot        uint32
	fractalGradientTextureID uint32
}

func createColorGradient(maxIterations int32, saturation float32, luminosity float32) []mgl32.Vec3 {
	colorGradient := make([]mgl32.Vec3, maxIterations)

	for idx := range colorGradient {
		hue := float32(idx) / (float32(maxIterations) * 1.15) // We offset this a bit so we end solidly in purple and not back in red.
		color := color.HSL{H: float64(hue), S: float64(saturation), L: float64(luminosity)}.ToRGB()
		colorGradient[idx] = mgl32.Vec3{float32(color.R), float32(color.G), float32(color.B)}
	}

	return colorGradient
}

func NewRainbowTexture(maxIterations int32) *RainbowTexture {
	var rainbow RainbowTexture

	colorGradient := createColorGradient(maxIterations, 1.0, 0.5)

	rainbow.openGlTextureSlot = gl.TEXTURE0
	gl.ActiveTexture(rainbow.openGlTextureSlot)
	gl.GenTextures(1, &rainbow.fractalGradientTextureID)

	gl.BindTexture(gl.TEXTURE_1D, rainbow.fractalGradientTextureID)
	gl.TexStorage1D(gl.TEXTURE_1D, 1, gl.RGB32F, maxIterations)
	gl.TexSubImage1D(gl.TEXTURE_1D, 0, 0, maxIterations, gl.RGB, gl.FLOAT, gl.Ptr(colorGradient))

	return &rainbow
}

func (rainbow *RainbowTexture) Delete() {
	gl.DeleteTextures(1, &rainbow.fractalGradientTextureID)
}
