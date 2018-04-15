package text

import (
	"fmt"
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/utils"
	"image"
	"image/draw"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var pixelsToVerticesScale float32 = 0.05 // Scales down the pixel size of a character to vertices
var runeFontSize int = 72                // Large enough to not look pixellated, small enough to be reasonable.

// Defines the index of a character in the texture maps
type characterIndex struct {
	CharacterOffset mgl32.Vec2 // Offset (already scaled) when drawing the character from the baseline

	MinBounds     utils.IntVec2 // Bounds of the character *in pixels.*
	MaxBounds     utils.IntVec2
	FontTextureId uint32
}

func (c *characterIndex) MinBoundsAsTextureCoords(textureSize float32) mgl32.Vec2 {
	return mgl32.Vec2{float32(c.MinBounds.X()) / float32(textureSize), float32(c.MinBounds.Y()) / float32(textureSize)}
}

func (c *characterIndex) MaxBoundsAsTextureCoords(textureSize float32) mgl32.Vec2 {
	return mgl32.Vec2{float32(c.MaxBounds.X()) / float32(textureSize), float32(c.MaxBounds.Y()) / float32(textureSize)}
}

type TextRenderer struct {
	context *freetype.Context
	font    *truetype.Font

	program textRendererProgram
	buffers textProgramBuffers

	halfMaxTextureSize int32
	fontTextures       []uint32
	nextLineOffset     int
	currentOffset      utils.IntVec2

	// Given a character, returns where it is on the textures for drawing
	characterMap map[rune]characterIndex
}

func (r *TextRenderer) preRender() {
	gl.UseProgram(r.program.shaderProgram)
	gl.BindVertexArray(r.buffers.vao)
}

// Renders the given rune using the provided model matrix and text-based offset.
// preRender(...) must be called before this method is called.
// Returns the x-offset of the text that was drawn.
func (r *TextRenderer) render(character rune, offset float32, model *mgl32.Mat4) float32 {
	// TODO: Add or get rune, position appropriately, and render, returning the character information.
	runeData := r.addOrGetRuneData(character)

	gl.ActiveTexture(gl.TEXTURE0 + runeData.FontTextureId)
	gl.BindTexture(gl.TEXTURE_2D, r.fontTextures[runeData.FontTextureId])
	gl.Uniform1i(r.program.fontImageLoc, int32(runeData.FontTextureId))

	gl.UniformMatrix4fv(r.program.modelLoc, 1, false, &model[0])

	runeOffset := runeData.CharacterOffset.Add(mgl32.Vec2{offset, 0})
	sendPrimitivesToDevice(r.buffers.positionVbo, r.buffers.colorVbo, r.buffers.texPosVbo,
		runeOffset,
		runeData.MinBoundsAsTextureCoords(float32(r.halfMaxTextureSize)),
		runeData.MaxBoundsAsTextureCoords(float32(r.halfMaxTextureSize)))
	renderPrimitive()

	return float32(runeData.MaxBounds.X()-runeData.MinBounds.X())*pixelsToVerticesScale + runeData.CharacterOffset.X()
}

func (r *TextRenderer) Delete() {
	r.buffers.Delete()
	r.program.Delete()
}

func loadContext(fontFileName string) (*truetype.Font, *freetype.Context) {
	fontFile := utils.ReadFileAsBytes(fontFileName)

	// Loads all the ASCII printable characters
	context := freetype.NewContext()
	parsedFont, err := freetype.ParseFont(fontFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse a TrueType font from the font file: %v", err))
	}

	context.SetDPI(72.0)
	context.SetFontSize(float64(runeFontSize))
	context.SetHinting(font.HintingFull)
	context.SetFont(parsedFont)

	return parsedFont, context
}

// Adds a rune to the list of characters
func (renderer *TextRenderer) addRune(character rune) {
	runeIndex := renderer.font.Index(character)
	fixedRuneFontSize := fixed.I(runeFontSize)

	hMetric := renderer.font.HMetric(fixedRuneFontSize, runeIndex)
	vMetric := renderer.font.VMetric(fixedRuneFontSize, runeIndex)

	bounds := renderer.font.Bounds(fixedRuneFontSize)
	boundRanges := bounds.Max.Sub(bounds.Min)

	// Compute image height so we just draw the character inside the box.
	maxHeight := boundRanges.Y.Ceil() - vMetric.TopSideBearing.Ceil()
	maxWidth := boundRanges.X.Ceil() - hMetric.LeftSideBearing.Ceil()

	// Move to the next line if needed
	if maxWidth+renderer.currentOffset.X() >= int(renderer.halfMaxTextureSize) {
		renderer.currentOffset[0] = 0

		if renderer.nextLineOffset == -1 {
			panic("Attempted to draw a character that is wider than the texture image. This is not supported.")
		}

		renderer.currentOffset[1] += renderer.nextLineOffset
	}

	// We have filled this texture image, so move onto the next one.
	if maxHeight+renderer.currentOffset.Y() >= int(renderer.halfMaxTextureSize) {
		renderer.addFontTexture()
		renderer.currentOffset = utils.IntVec2{0, 0}
		renderer.nextLineOffset = -1
	}

	dstImage := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))
	draw.Draw(dstImage, dstImage.Bounds(), image.White, image.ZP, draw.Src)

	renderer.context.SetClip(dstImage.Bounds())
	renderer.context.SetSrc(image.Black)
	renderer.context.SetDst(dstImage)

	xOffset := -hMetric.LeftSideBearing.Ceil()
	yHeight := vMetric.TopSideBearing.Ceil()

	// Draw, copy, and save the new character.
	point, err := renderer.context.DrawString(string(character), freetype.Pt(xOffset, yHeight))
	if err != nil {
		panic(fmt.Sprintf("Unable to draw rune '%v' to destination %v : %v", character, point, err))
	}

	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		int32(renderer.currentOffset.X()), int32(renderer.currentOffset.Y()),
		int32(maxWidth), int32(maxHeight),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dstImage.Pix))

	renderer.characterMap[character] = characterIndex{
		FontTextureId: uint32(len(renderer.fontTextures) - 1),
		MinBounds:     renderer.currentOffset,
		MaxBounds:     utils.IntVec2{renderer.currentOffset.X() + maxWidth, renderer.currentOffset.Y() + maxHeight},
		CharacterOffset: mgl32.Vec2{
			float32(hMetric.LeftSideBearing.Ceil()) * pixelsToVerticesScale,
			float32(vMetric.TopSideBearing.Ceil()) * pixelsToVerticesScale,
		}}

	// Update the offset and save our rune
	renderer.currentOffset[0] += maxWidth
	if maxHeight > renderer.nextLineOffset {
		renderer.nextLineOffset = maxHeight
	}
}

func (r *TextRenderer) addOrGetRuneData(runeChar rune) characterIndex {
	if _, hasRune := r.characterMap[runeChar]; !hasRune {
		r.addRune(runeChar)
	}

	return r.characterMap[runeChar]
}

// Adds in a new font texture
func (r *TextRenderer) addFontTexture() {
	maxTextures := opengl.GetGlCaps().MaxTextures
	if int32(len(r.fontTextures)) >= maxTextures {
		howToFix := "Either reduce the number of unique characters being rendered or upgrade your graphics hardware."
		panic(fmt.Sprintf("Cannot add a new font texture as we've exceeded the maximum number of textures (%v).\n%v\n", maxTextures, howToFix))
	}

	var newTextureId uint32
	gl.GenTextures(1, &newTextureId)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(len(r.fontTextures)))
	gl.BindTexture(gl.TEXTURE_2D, newTextureId)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, r.halfMaxTextureSize, r.halfMaxTextureSize)

	r.fontTextures = append(r.fontTextures, newTextureId)
}

func NewTextRenderer(fontFile string) *TextRenderer {
	renderer := TextRenderer{
		nextLineOffset: -1,
		currentOffset:  utils.IntVec2{0, 0},
		fontTextures:   make([]uint32, 0),
		characterMap:   make(map[rune]characterIndex)}

	renderer.program = newTextRendererProgram()
	renderer.buffers = newTextProgramBuffers()
	renderer.font, renderer.context = loadContext(fontFile)

	renderer.halfMaxTextureSize = opengl.GetGlCaps().MaxTextureSize
	renderer.addFontTexture()

	return &renderer
}

// Retrieves the size that a given string will be when rendered without any model manipulation.
func (r *TextRenderer) GetSize(text string) mgl32.Vec2 {
	size := mgl32.Vec2{0, 0}

	// TODO: we probably need to actually return an array so rendering functions can take care of alignment themselves.
	for _, char := range text {
		characterInfo := r.addOrGetRuneData(char)
		size[0] += float32(characterInfo.MaxBounds.X()-characterInfo.MinBounds.X())*pixelsToVerticesScale + characterInfo.CharacterOffset.X()
		size[1] = utils.MaxFloat32(size[1], float32(characterInfo.MaxBounds.Y())*pixelsToVerticesScale+characterInfo.CharacterOffset.Y())
	}

	return size
}

// Implement Renderer
func (r *TextRenderer) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(r.program.shaderProgram)
	gl.UniformMatrix4fv(r.program.projectionLoc, 1, false, &projection[0])
}

func (r *TextRenderer) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(r.program.shaderProgram)
	gl.UniformMatrix4fv(r.program.cameraLoc, 1, false, &camera[0])
}
