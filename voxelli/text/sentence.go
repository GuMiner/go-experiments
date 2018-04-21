package text

import (
	"go-experiments/voxelli/utils"

	"github.com/go-gl/mathgl/mgl32"
)

type Sentence struct {
	TextRenderer *TextRenderer
	Background   mgl32.Vec3
	Foreground   mgl32.Vec3
}

// Renders the given text character by character with minimal overhead.
// Useful for small snippets of text that change frequently
func (r *Sentence) GetRenderSize(text string) mgl32.Vec2 {
	aggregateSize := mgl32.Vec2{0, 0}
	for _, runeChar := range text {
		size := r.TextRenderer.getCharacterSize(runeChar)
		aggregateSize[0] += size[0]
		aggregateSize[1] = utils.MaxFloat32(aggregateSize[1], size[1])
	}

	return aggregateSize
}

func (r *Sentence) Render(text string, model *mgl32.Mat4, doubleSided bool) {
	r.TextRenderer.preRender(r.Background, r.Foreground, model)

	currentOffset := float32(0.0)
	for _, runeChar := range text {
		characterOffset := r.TextRenderer.render(runeChar, currentOffset)
		currentOffset += characterOffset
	}

	// Do the same thing in reverse now, so that the text is visible from both sides
	if doubleSided {
		reverseOffset := currentOffset
		currentOffset = 0.0
		for _, runeChar := range text {
			characterOffset := r.TextRenderer.renderReverse(runeChar, currentOffset, reverseOffset)
			currentOffset += characterOffset
		}
	}
}

func NewSentence(renderer *TextRenderer, background, foreground mgl32.Vec3) *Sentence {
	sentence := Sentence{TextRenderer: renderer, Background: background, Foreground: foreground}
	return &sentence
}
