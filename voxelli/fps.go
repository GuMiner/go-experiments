// Defines a simple FPS counter
package main

import (
	"fmt"
	"go-experiments/voxelli/text"

	"github.com/go-gl/mathgl/mgl32"
)

type FpsCounter struct {
	// Settings
	refreshInterval float32
	sentence        *text.Sentence
	offsetPosition  mgl32.Vec3

	frameCount  int
	elapsedTime float32
	currentFps  float32
}

func NewFpsCounter(fpsSentence *text.Sentence, refreshInterval float32, offsetPosition mgl32.Vec3) *FpsCounter {
	fpsCounter := FpsCounter{
		refreshInterval: refreshInterval,
		sentence:        fpsSentence,
		offsetPosition:  offsetPosition,
		frameCount:      0,
		elapsedTime:     0,
		currentFps:      0.0}

	return &fpsCounter
}

func (f *FpsCounter) Update(frameTime float32) {
	f.elapsedTime += frameTime
	f.frameCount++

	if f.elapsedTime > f.refreshInterval {
		f.currentFps = float32(f.frameCount) / f.elapsedTime

		f.elapsedTime = 0.0
		f.frameCount = 0
	}
}

func InverseRotationMatrix(eye, center, up mgl32.Vec3) mgl32.Mat4 {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)

	M := mgl32.Mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		0, 0, 0, 1,
	}

	return M.Transpose() // Inverse of a rotation matrix
}

func (f *FpsCounter) Render(camera *Camera) {
	fpsString := fmt.Sprintf("FPS: %.2f", f.currentFps)

	// Put the text in front of the camera, scaled appropriately.
	objectPos := camera.Position.Add(camera.Forwards.Mul(1))
	frontOfCamera := mgl32.Translate3D(objectPos.X(), objectPos.Y(), objectPos.Z())

	offset := mgl32.Translate3D(f.offsetPosition.X(), f.offsetPosition.Y(), 0)
	scale := mgl32.Scale3D(f.offsetPosition.Z(), f.offsetPosition.Z(), 1)

	// Get the text so we're rotating from the center and the front is in the Y+ direction.
	right := mgl32.Vec3{1, 0, 0}
	tiltRotation := mgl32.HomogRotate3D(0, right)

	renderSize := f.sentence.GetRenderSize(fpsString)
	center := mgl32.Translate3D(-renderSize.X()/2, -renderSize.Y()/2, 0)

	// Reverse the camera rotation
	orientRotation := InverseRotationMatrix(camera.Position, camera.Target, camera.Up)

	fpsModelMatrix := frontOfCamera.Mul4(
		orientRotation.Mul4(
			offset.Mul4(
				scale.Mul4(
					tiltRotation.Mul4(center)))))
	f.sentence.Render(fpsString, &fpsModelMatrix, true)
}
