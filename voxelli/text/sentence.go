package text

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Sentence struct {
	textRenderer *TextRenderer
	Background   mgl32.Vec3
	Foreground   mgl32.Vec3
}

type SentenceBatch struct {
}

type OptimizedSentence struct {
	vao         uint32
	positionVbo uint32
	colorVbo    uint32
	texPosVbo   uint32
	sentence    Sentence
}

// Renders the given text character by character with minimal overhead.
// Useful for small snippets of text that change frequently
func (r *Sentence) Render(text string, model *mgl32.Mat4) {
	r.textRenderer.preRender()

	for _, runeChar := range text {
		r.textRenderer.render(runeChar, model)
	}
}

// Batches up the text to render all characters at once.
// Useful for large snippets of text that change.
func (r *Sentence) Batch(text string) *SentenceBatch {
	var batch SentenceBatch
	return &batch
}

// Creates a custom texture that contains the entire contents of the given text.
// Useful for text that doesn't change.
func (r *Sentence) Optimize(text string) *OptimizedSentence {
	var optimized OptimizedSentence
	return &optimized
}

// Renders a pre-batched set of text
func (r *SentenceBatch) Render(model *mgl32.Mat4) {
}

// Deletes resources allocated by the SentenceBatch
func (r *SentenceBatch) Delete() {
}

// Renders an optimized set of text
func (r *OptimizedSentence) Render(model *mgl32.Mat4) {
}

// Deletes resources allocated by the OptimizedSentence
func (r *OptimizedSentence) Delete() {
}

func NewSentence(renderer *TextRenderer) *Sentence {
	sentence := Sentence{textRenderer: renderer}
	return &sentence
}
