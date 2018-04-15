package vehicle

import (
	"go-experiments/voxelli/input"
)

var isColorOverlayEnabled bool = true

// Checks if the wireframe button has been toggled or not, toggling the GL setting
// This function should be called within the OpenGL update loop
func CheckColorOverlayToggle() {
	if input.IsTyped(input.ToggleColorOverlay) {
		isColorOverlayEnabled = !isColorOverlayEnabled
	}
}
