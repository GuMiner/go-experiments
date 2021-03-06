package text

import (
	"fmt"
	"go-experiments/common/commonconfig"
	"go-experiments/common/commonio"
	"go-experiments/common/commonmath"
	"go-experiments/common/commonopengl"
	"go-experiments/voxelli/text/renderer"
	"image"
	"image/draw"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

// Defines the index of a character in the texture maps
type characterIndex struct {
	Offset        commonMath.IntVec2 // Bounds of the character *in pixels.*
	Scale         commonMath.IntVec2
	FontTextureId uint32
}

type TextRenderer struct {
	Program *textRenderer.TextRendererProgram

	context *freetype.Context
	font    *truetype.Font

	buffers *textRenderer.TextProgramBuffers

	textureSize    int32
	fontTextures   []uint32
	nextLineOffset int
	currentOffset  commonMath.IntVec2

	// Given a character, returns where it is on the textures for drawing
	characterMap map[rune]characterIndex
}

func (r *TextRenderer) preRender(background, foreground mgl32.Vec3, model *mgl32.Mat4) {
	r.Program.UseProgram(r.buffers)
	r.Program.SetColors(background, foreground)
	r.Program.SetModel(model)
}

// Renders the given rune using the provided model matrix and text-based offset.
// preRender(...) must be called before this method is called.
// Returns the x-offset of the text that was drawn.
func (r *TextRenderer) render(character rune, offset float32) float32 {
	runeData := r.addOrGetRuneData(character)

	r.Program.SetTexture(runeData.FontTextureId, r.fontTextures[runeData.FontTextureId])
	positionBuffer, uvBuffer, runeOffset := generateCharacterPrimitive(
		offset,
		runeData.Offset, runeData.Scale,
		r.textureSize,
		false)

	r.buffers.SendToDevice(positionBuffer, uvBuffer)
	renderPrimitive()

	return runeOffset
}

func (r *TextRenderer) getCharacterSize(character rune) mgl32.Vec2 {
	runeData := r.addOrGetRuneData(character)
	xScale, yScale := computeCharacterScale(runeData.Scale)
	return mgl32.Vec2{xScale, yScale}
}

// Same as 'render' but starts at an offset and flips the characters around, for double-sided displays
func (r *TextRenderer) renderReverse(character rune, offset float32, reverseOffset float32) float32 {
	runeData := r.addOrGetRuneData(character)

	r.Program.SetTexture(runeData.FontTextureId, r.fontTextures[runeData.FontTextureId])

	characterXScale, _ := computeCharacterScale(runeData.Scale)
	positionBuffer, uvBuffer, runeOffset := generateCharacterPrimitive(
		reverseOffset-(offset+characterXScale),
		runeData.Offset, runeData.Scale,
		r.textureSize,
		true)

	r.buffers.SendToDevice(positionBuffer, uvBuffer)
	renderPrimitive()

	return runeOffset
}

func (r *TextRenderer) Delete() {
	r.buffers.Delete()
	r.Program.Delete()
}

func loadContext(fontFileName string) (*truetype.Font, *freetype.Context) {
	fontFile := commonIo.ReadFileAsBytes(fontFileName)

	// Loads all the ASCII printable characters
	context := freetype.NewContext()
	parsedFont, err := freetype.ParseFont(fontFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse a TrueType font from the font file: %v", err))
	}

	context.SetDPI(float64(commonConfig.Config.Text.RuneFontSize))
	context.SetFontSize(float64(commonConfig.Config.Text.RuneFontSize))
	context.SetHinting(font.HintingFull)
	context.SetFont(parsedFont)

	return parsedFont, context
}

// Advances to the next line or image in the set of texture images, as neessary
func (renderer *TextRenderer) advanceIfNecessary(width, height int) {
	// Move to the next line if needed
	if width+renderer.currentOffset.X() >= int(renderer.textureSize) {
		renderer.currentOffset[0] = 0

		if renderer.nextLineOffset == -1 {
			panic("Attempted to draw a character that is wider than the texture image. This is not supported.")
		}

		renderer.currentOffset[1] += renderer.nextLineOffset
		renderer.nextLineOffset = -1
	}

	// We have filled this texture image, so move onto the next one.
	if height+renderer.currentOffset.Y() >= int(renderer.textureSize) {
		renderer.addFontTexture()
		renderer.currentOffset = commonMath.IntVec2{0, 0}
		renderer.nextLineOffset = -1
	}
}

func (renderer *TextRenderer) updateRuneOffset(width, height int) {
	// Update the offset and save our rune
	renderer.currentOffset[0] += width
	if height > renderer.nextLineOffset {
		renderer.nextLineOffset = height
	}
}

// Adds a rune to the list of characters
func (renderer *TextRenderer) addRune(character rune) {
	runeIndex := renderer.font.Index(character)
	fixedRuneFontSize := fixed.I(commonConfig.Config.Text.RuneFontSize)

	hMetric := renderer.font.HMetric(fixedRuneFontSize, runeIndex)
	vMetric := renderer.font.VMetric(fixedRuneFontSize, runeIndex)

	// Compute image height so we just draw the character inside the box.
	maxWidth := hMetric.AdvanceWidth.Ceil() - hMetric.LeftSideBearing.Ceil()
	maxHeight := vMetric.AdvanceHeight.Ceil()

	fullWidth := maxWidth + commonConfig.Config.Text.BorderSize*2
	fullHeight := maxHeight + commonConfig.Config.Text.BorderSize*2
	renderer.advanceIfNecessary(fullWidth, fullHeight)

	dstImage := image.NewRGBA(image.Rect(0, 0, fullWidth, fullHeight))
	draw.Draw(dstImage, dstImage.Bounds(), image.White, image.ZP, draw.Src)

	renderer.context.SetClip(dstImage.Bounds())
	renderer.context.SetSrc(image.Black)
	renderer.context.SetDst(dstImage)

	xOffset := -hMetric.LeftSideBearing.Ceil() + commonConfig.Config.Text.BorderSize
	yHeight := vMetric.TopSideBearing.Ceil() + commonConfig.Config.Text.BorderSize

	// Draw, copy, and save the new character.
	point, err := renderer.context.DrawString(string(character), freetype.Pt(xOffset, yHeight))
	if err != nil {
		panic(fmt.Sprintf("Unable to draw rune '%v' to destination %v : %v", character, point, err))
	}

	gl.BindTexture(gl.TEXTURE_2D, renderer.fontTextures[len(renderer.fontTextures)-1])
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		int32(renderer.currentOffset.X()), int32(renderer.currentOffset.Y()),
		int32(fullWidth), int32(fullHeight),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dstImage.Pix))

	renderer.characterMap[character] = characterIndex{
		FontTextureId: uint32(len(renderer.fontTextures) - 1),
		Offset:        renderer.currentOffset,
		Scale:         commonMath.IntVec2{fullWidth, fullHeight}}

	renderer.updateRuneOffset(fullWidth, fullHeight)
}

func (r *TextRenderer) addOrGetRuneData(runeChar rune) characterIndex {
	if _, hasRune := r.characterMap[runeChar]; !hasRune {
		r.addRune(runeChar)
	}

	return r.characterMap[runeChar]
}

// Adds in a new font texture
func (r *TextRenderer) addFontTexture() {
	maxTextures := commonOpenGl.GetGlCaps().MaxTextures
	if int32(len(r.fontTextures)) >= maxTextures {
		howToFix := "Either reduce the number of unique characters being rendered or upgrade your graphics hardware."
		panic(fmt.Sprintf("Cannot add a new font texture as we've exceeded the maximum number of textures (%v).\n%v\n", maxTextures, howToFix))
	}

	var newTextureId uint32
	gl.GenTextures(1, &newTextureId)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(len(r.fontTextures)))
	gl.BindTexture(gl.TEXTURE_2D, newTextureId)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, r.textureSize, r.textureSize)

	r.fontTextures = append(r.fontTextures, newTextureId)
}

func NewTextRenderer(fontFile string) *TextRenderer {
	renderer := TextRenderer{
		nextLineOffset: -1,
		currentOffset:  commonMath.IntVec2{0, 0},
		fontTextures:   make([]uint32, 0),
		characterMap:   make(map[rune]characterIndex)}

	renderer.Program = textRenderer.NewTextRendererProgram()
	renderer.buffers = textRenderer.NewTextProgramBuffers()
	renderer.font, renderer.context = loadContext(fontFile)

	renderer.textureSize = commonMath.MinInt32(
		commonOpenGl.GetGlCaps().MaxTextureSize, 2048)
	renderer.addFontTexture()

	return &renderer
}

// Implement Renderer
func (r *TextRenderer) UpdateProjection(projection *mgl32.Mat4) {
	r.Program.UpdateProjection(projection)
}

func (r *TextRenderer) UpdateCamera(camera *mgl32.Mat4) {
	r.Program.UpdateCamera(camera)
}
