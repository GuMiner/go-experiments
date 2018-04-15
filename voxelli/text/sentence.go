package text

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Sentence struct {
	textRenderer *TextRenderer
	Background   mgl32.Vec3
	Foreground   mgl32.Vec3
}

// Renders the given text character by character with minimal overhead.
// Useful for small snippets of text that change frequently
func (r *Sentence) Render(text string, model *mgl32.Mat4) {
	r.textRenderer.preRender()

	// TODO: Use background and foreground, by sending to shader.
	currentOffset := float32(0.0)
	for _, runeChar := range text {
		characterOffset := r.textRenderer.render(runeChar, currentOffset, model)
		currentOffset += characterOffset
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

func NewSentence(renderer *TextRenderer, background, foreground mgl32.Vec3) *Sentence {
	sentence := Sentence{textRenderer: renderer, Background: background, Foreground: foreground}
	return &sentence
}
