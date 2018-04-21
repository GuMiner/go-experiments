package main

import (
	"fmt"
	"go-experiments/voxelli/text"
	"go-experiments/voxelli/utils"

	"github.com/go-gl/mathgl/mgl32"
)

type HelpText struct {
	sentence       *text.Sentence
	offsetPosition mgl32.Vec3
}

func NewHelpText(helpTextSentence *text.Sentence, offsetPosition mgl32.Vec3) *HelpText {
	helpText := HelpText{
		sentence:       helpTextSentence,
		offsetPosition: offsetPosition}

	return &helpText
}

func (f *HelpText) Render(camera *Camera) {
	// TODO: Somehow get human-friendly versions of these, instead of hardcoding
	helpTextLines := []string{
		fmt.Sprintf("Display:"),
		fmt.Sprintf("  Color Overlay: C"),
		fmt.Sprintf("  Score Overlay: T"),
		fmt.Sprintf("  FPS: F"),
		fmt.Sprintf("  Help: H"),
		fmt.Sprintf("  Debug: G"),
		fmt.Sprintf("Motion:"),
		fmt.Sprintf("  Look Left: Left"),
		fmt.Sprintf("  Look Right: Right"),
		fmt.Sprintf("  Look Up: Up"),
		fmt.Sprintf("  Look Down: Down"),
		fmt.Sprintf("  Rotate CW: E"),
		fmt.Sprintf("  Rotate CCW: D"),
		fmt.Sprintf("  Move Forwards: A"),
		fmt.Sprintf("  Move Backwards: Z"),
		fmt.Sprintf("  Move Left: Q"),
		fmt.Sprintf("  Move Right: W"),
		fmt.Sprintf("  Move Up: S"),
		fmt.Sprintf("  Move Down: X")}

	// Put the text in front of the camera, scaled appropriately.
	objectPos := camera.Position.Add(camera.Forwards.Mul(1))
	frontOfCamera := mgl32.Translate3D(objectPos.X(), objectPos.Y(), objectPos.Z())

	scale := mgl32.Scale3D(f.offsetPosition.Z(), f.offsetPosition.Z(), 1)

	// Get the text so we're rotating from the center and the front is in the Y+ direction.
	right := mgl32.Vec3{1, 0, 0}
	tiltRotation := mgl32.HomogRotate3D(0, right)

	// Reverse the camera rotation
	orientRotation := utils.InverseLookAtRotationMatrix(camera.Position, camera.Target, camera.Up)

	for i, helpTextLine := range helpTextLines {

		renderSize := f.sentence.GetRenderSize(helpTextLine)
		center := mgl32.Translate3D(0, -renderSize.Y()/2, 0)

		offset := mgl32.Translate3D(f.offsetPosition.X(), f.offsetPosition.Y()-float32(i)/28.0, 0)

		fpsModelMatrix := frontOfCamera.Mul4(
			orientRotation.Mul4(
				offset.Mul4(
					scale.Mul4(
						tiltRotation.Mul4(center)))))
		f.sentence.Render(helpTextLine, &fpsModelMatrix, true)
	}
}
