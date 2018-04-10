package diagnostics

import (
	"go-experiments/voxelli/input"
)

var isDebug bool = false

func CheckDebugToggle() {
	if input.IsTyped(input.ToggleDebug) {
		isDebug = !isDebug
	}
}

func IsDebug() bool {
	return isDebug
}
