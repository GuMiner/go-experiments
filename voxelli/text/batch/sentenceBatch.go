package batch

import (
	"go-experiments/voxelli/text"
	"go-experiments/voxelli/text/renderer"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type BufferTexturePair struct {
	buffer          *textRenderer.TextProgramBuffers
	textureId       uint32
	textureResource uint32

	vertexCount int32
}

type SentenceBatch struct {
	sentence *text.Sentence
	buffers  []BufferTexturePair
}

// Batches up the text to render all characters at once.
// Useful for large snippets of text that change infrequently
func NewBatch(sentence *text.Sentence, text string, doubleSided bool) *SentenceBatch {
	batch := SentenceBatch{
		sentence: sentence,
		buffers:  make([]BufferTexturePair, 0)}

	return &batch
}

// Renders a pre-batched set of text
func (r *SentenceBatch) Render(model *mgl32.Mat4) {
	for _, bufferSet := range r.buffers {
		r.sentence.TextRenderer.Program.UseProgram(bufferSet.buffer)
		r.sentence.TextRenderer.Program.SetColors(r.sentence.Background, r.sentence.Foreground)
		r.sentence.TextRenderer.Program.SetModel(model)

		r.sentence.TextRenderer.Program.SetTexture(bufferSet.textureId, bufferSet.textureResource)

		gl.DrawArrays(gl.TRIANGLES, 0, bufferSet.vertexCount)
	}
}

// Deletes resources allocated by the SentenceBatch
func (r *SentenceBatch) Delete() {
	for _, bufferSet := range r.buffers {
		bufferSet.buffer.Delete()
	}
}
