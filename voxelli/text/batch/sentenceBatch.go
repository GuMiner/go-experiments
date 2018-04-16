package batch

import (
	"go-experiments/voxelli/text"

	"github.com/go-gl/mathgl/mgl32"
)

type SentenceBatch struct {
	sentence *text.Sentence
}

// Batches up the text to render all characters at once.
// Useful for large snippets of text that change infrequently
func NewBatch(sentence *text.Sentence, text string, doubleSided bool) *SentenceBatch {
	batch := SentenceBatch{sentence: sentence}
	return &batch
}

// Renders a pre-batched set of text
func (r *SentenceBatch) Render(model *mgl32.Mat4) {

}

// Deletes resources allocated by the SentenceBatch
func (r *SentenceBatch) Delete() {
}
