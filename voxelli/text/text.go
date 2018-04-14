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

// Defines the index of a character in the texture maps
type characterIndex struct {
	CharacterOffset mgl32.Vec2 // Offset (already scaled) when drawing the character from the baseline

	MinBounds     utils.IntVec2 // Bounds of the character in pixels.
	MaxBounds     utils.IntVec2
	FontTextureId uint32
}

type TextRenderer struct {
	context *freetype.Context
	font    *truetype.Font

	program textRendererProgram
	buffers textProgramBuffers

	halfMaxTextureSize int32
	fontTextures       []uint32
	nextLineOffset     int32
	currentOffset      utils.IntVec2

	// Given a character, returns where it is on the textures for drawing
	characterMap map[rune]characterIndex
}

func (r *TextRenderer) preRender() {
	gl.UseProgram(r.program.shaderProgram)
	gl.BindVertexArray(r.buffers.vao)
}

// Renders the given rune using the provided model matrix.
// preRender(...) must be called before this method is called.
func (r *TextRenderer) render(character rune, model *mgl32.Mat4) mgl32.Vec2 {
	// TODO: Add or get rune, position appropriately, and render, returning the character information.
	runeData := r.addOrGetRuneData(character)

	gl.ActiveTexture(gl.TEXTURE0 + runeData.FontTextureId)
	gl.BindTexture(gl.TEXTURE_2D, r.fontTextures[runeData.FontTextureId])
	gl.Uniform1i(r.program.fontImageLoc, int32(runeData.FontTextureId))

	// TODO: First, this should be in character primitive. Second, we need to pass in the UV coordinates and the color foreground / background data.
	// Third, that means we need to pass in that data here.
	gl.UniformMatrix4fv(r.program.modelLoc, 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	return mgl32.Vec2{-1, -1}
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
	context.SetFontSize(72.0) // Large enough to not look pixellated, small enough to be reasonable.
	context.SetHinting(font.HintingFull)
	context.SetFont(parsedFont)

	return parsedFont, context
}

// Adds a rune to the list of characters
func (renderer *TextRenderer) addRune(character rune) {
	dstImage := image.NewRGBA(image.Rect(0, 0, int(renderer.halfMaxTextureSize), int(renderer.halfMaxTextureSize)))
	draw.Draw(dstImage, dstImage.Bounds(), image.White, image.ZP, draw.Src)

	renderer.context.SetClip(dstImage.Bounds())
	renderer.context.SetSrc(image.Black)
	renderer.context.SetDst(dstImage)

	stringToDraw := "J23456789 abcdefg!-=?αΩ♣"

	// Get max vertical ascent
	maxSideBearing := 0
	for _, runeVal := range stringToDraw {
		vmetric := renderer.font.VMetric(fixed.I(16), renderer.font.Index(runeVal))
		sideBearing := vmetric.TopSideBearing.Ceil()
		if sideBearing > maxSideBearing {
			maxSideBearing = sideBearing
		}
	}
	//
	// indexTwo := renderer.font.Index('2')
	//
	hmetric := renderer.font.HMetric(fixed.I(16), renderer.font.Index('J')) // If negative, we need to offset the string by this amount.
	fmt.Printf("%v, %v\n\n", hmetric.AdvanceWidth.Floor(), hmetric.LeftSideBearing.Floor())
	//
	// kern := renderer.font.Kern(fixed.I(16), indexOne, indexTwo)
	// fmt.Printf("%v\n", kern.Floor())

	bounds := renderer.font.Bounds(fixed.I(16))
	boundRanges := bounds.Max.Sub(bounds.Min)

	// Offset height
	yHeight := boundRanges.Y.Ceil() - maxSideBearing

	// index := renderer.font.Index('J')

	point, err := renderer.context.DrawString(stringToDraw, freetype.Pt(0, yHeight))
	if err != nil {
		panic(fmt.Sprintf("Unable to draw string to destination %v : %v", point, err))
	}

	// fmt.Printf("%v", dstImage.Pix)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, 512, 512, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dstImage.Pix))
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
		nextLineOffset: 0,
		currentOffset:  utils.IntVec2{0, 0}}

	renderer.program = newTextRendererProgram()
	renderer.font, renderer.context = loadContext(fontFile)

	renderer.halfMaxTextureSize = opengl.GetGlCaps().MaxTextureSize
	renderer.fontTextures = make([]uint32, 0)
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
