package text

import (
	"github.com/go-gl/mathgl/mgl32"
)

type OptimizedSentence struct {
	vao         uint32
	positionVbo uint32
	colorVbo    uint32
	texPosVbo   uint32
	sentence    Sentence
}

// Renders an optimized set of text
func (r *OptimizedSentence) Render(model *mgl32.Mat4) {
}

// Deletes resources allocated by the OptimizedSentence
func (r *OptimizedSentence) Delete() {
}
