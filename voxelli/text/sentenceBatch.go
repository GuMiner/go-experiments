package text

import (
	"github.com/go-gl/mathgl/mgl32"
)

type SentenceBatch struct {
	sentence Sentence
}

// Renders a pre-batched set of text
func (r *SentenceBatch) Render(model *mgl32.Mat4) {
}

// Deletes resources allocated by the SentenceBatch
func (r *SentenceBatch) Delete() {
}
