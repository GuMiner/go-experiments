package text

import (
	"fmt"
	"go-experiments/voxelli/opengl"
	"go-experiments/voxelli/utils"
	"image"
	"image/draw"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

type Sentence struct {
	vao         uint32
	positionVbo uint32
	colorVbo    uint32
	texPosVbo   uint32
}

// Temporary to demo drawing
var ccwQuadVert = []mgl32.Vec3{
	mgl32.Vec3{-50, 50, 40},
	mgl32.Vec3{50, -50, 40},
	mgl32.Vec3{50, 50, 40},

	mgl32.Vec3{-50, -50, 40},
	mgl32.Vec3{50, -50, 40},
	mgl32.Vec3{-50, 50, 40}}

var ccwQuadColor = []mgl32.Vec3{
	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5},

	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5},
	mgl32.Vec3{1.0, 1.0, 0.5}}

var ccwQuadUv = []mgl32.Vec2{
	mgl32.Vec2{0, 0},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{1, 0},

	mgl32.Vec2{0, 1},
	mgl32.Vec2{1, 1},
	mgl32.Vec2{0, 0}}

func NewSentence() *Sentence {
	var sentence Sentence

	gl.GenVertexArrays(1, &sentence.vao)
	gl.BindVertexArray(sentence.vao)

	gl.EnableVertexAttribArray(0)
	gl.GenBuffers(1, &sentence.positionVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, sentence.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadVert)*3*4, gl.Ptr(ccwQuadVert), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(1)
	gl.GenBuffers(1, &sentence.colorVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, sentence.colorVbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadColor)*3*4, gl.Ptr(ccwQuadColor), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(2)
	gl.GenBuffers(1, &sentence.texPosVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, sentence.texPosVbo)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 0, nil)

	gl.BufferData(gl.ARRAY_BUFFER, len(ccwQuadUv)*3*4, gl.Ptr(ccwQuadUv), gl.STATIC_DRAW)

	return &sentence
}

func (s *Sentence) Delete() {
	gl.DeleteBuffers(1, &s.texPosVbo)
	gl.DeleteBuffers(1, &s.colorVbo)
	gl.DeleteBuffers(1, &s.positionVbo)
	gl.DeleteVertexArrays(1, &s.vao)
}

type TextRenderer struct {
	context *freetype.Context

	shaderProgram uint32
	projectionLoc int32
	cameraLoc     int32
	modelLoc      int32
	fontImageLoc  int32

	fontTexture uint32

	// TODO: have more than one after we validate rendering works
	sentence *Sentence
}

func (r *TextRenderer) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(r.shaderProgram)
	gl.UniformMatrix4fv(r.projectionLoc, 1, false, &projection[0])
}

func (r *TextRenderer) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(r.shaderProgram)
	gl.UniformMatrix4fv(r.cameraLoc, 1, false, &camera[0])
}

func (r *TextRenderer) Render(text string, model *mgl32.Mat4) {
	gl.UseProgram(r.shaderProgram)

	gl.BindVertexArray(r.sentence.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.fontTexture)
	gl.Uniform1i(r.fontImageLoc, 0)

	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (r *TextRenderer) Delete() {
	r.sentence.Delete()
	gl.DeleteProgram(r.shaderProgram)
}

func loadContext(fontFileName string) *freetype.Context {
	fontFile := utils.ReadFileAsBytes(fontFileName)

	// Loads all the ASCII printable characters
	context := freetype.NewContext()
	parsedFont, err := freetype.ParseFont(fontFile)
	if err != nil {
		panic("Failed to parse a TrueType font from the font file!")
	}

	context.SetDPI(72.0)
	context.SetFontSize(16.0)
	context.SetHinting(font.HintingFull)
	context.SetFont(parsedFont)

	return context
}

func NewTextRenderer(fontFile string) *TextRenderer {
	var renderer TextRenderer
	// Setup shader
	renderer.shaderProgram = opengl.CreateProgram("./shaders/textRenderer")

	renderer.projectionLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("projection\x00"))
	renderer.cameraLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("camera\x00"))
	renderer.modelLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("model\x00"))
	renderer.fontImageLoc = gl.GetUniformLocation(renderer.shaderProgram, gl.Str("fontImage\x00"))

	renderer.sentence = NewSentence()

	// Setup font
	renderer.context = loadContext(fontFile)

	gl.GenTextures(1, &renderer.fontTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, renderer.fontTexture)

	// TODO: Implement dynamic scaling as needed
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, 512, 512)

	dstImage := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(dstImage, dstImage.Bounds(), image.White, image.ZP, draw.Src)

	renderer.context.SetClip(dstImage.Bounds())
	renderer.context.SetSrc(image.Black)
	renderer.context.SetDst(dstImage)

	point, err := renderer.context.DrawString("123456789 abcdefg!-=?αΩ♣", freetype.Pt(10, 10))
	if err != nil {
		panic(fmt.Sprintf("Unable to draw string to destination %v : %v", point, err))
	}

	// fmt.Printf("%v", dstImage.Pix)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, 512, 512, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dstImage.Pix))
	return &renderer
}
