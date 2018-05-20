package lines

import (
	"go-experiments/common/commonopengl"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type LinesShaderProgram struct {
	program uint32

	colorLoc int32

	vao uint32
	vbo uint32
}

func (r *LinesShaderProgram) sendLineDataToShader(lines [][2]mgl32.Vec2) {
	gl.BindVertexArray(r.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)

	flattenedLineData := make([]mgl32.Vec2, len(lines)*2)
	for i := 0; i < len(lines); i++ {
		flattenedLineData[i*2] = lines[i][0]
		flattenedLineData[i*2+1] = lines[i][1]
	}

	// 2 -- 2 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(flattenedLineData)*2*4, gl.Ptr(flattenedLineData), gl.DYNAMIC_DRAW)
}

func NewLinesShaderProgram() *LinesShaderProgram {
	lineShaderProg := LinesShaderProgram{}

	// Setup the OpenGL program
	lineShaderProg.program = commonOpenGl.CreateProgram("./ui/shaders/lines")

	lineShaderProg.colorLoc = gl.GetUniformLocation(lineShaderProg.program, gl.Str("givenColor\x00"))

	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &lineShaderProg.vao)
	gl.BindVertexArray(lineShaderProg.vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &lineShaderProg.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, lineShaderProg.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	return &lineShaderProg
}

func (shaderProgram *LinesShaderProgram) PreRender() {
	gl.UseProgram(shaderProgram.program)
	gl.BindVertexArray(shaderProgram.vao)
}

func (shaderProgram *LinesShaderProgram) Render(lines [][2]mgl32.Vec2, color mgl32.Vec3) {
	gl.Uniform3f(shaderProgram.colorLoc, color.X(), color.Y(), color.Z())

	// Break lines down into sets
	lineIdx := 0
	lineStep := 256

	for lineIdx < len(lines) {
		var localStep int32
		if lineIdx+lineStep < len(lines) {
			localStep = int32(lineStep)
		} else {
			localStep = int32(len(lines) - lineIdx)
		}

		shaderProgram.sendLineDataToShader(lines[lineIdx : lineIdx+int(localStep)])
		gl.DrawArrays(gl.LINES, 0, localStep*2) // 2 vertices per line

		lineIdx += int(localStep)
	}
}

func (shaderProgram *LinesShaderProgram) Delete() {
	gl.DeleteBuffers(1, &shaderProgram.vbo)
	gl.DeleteVertexArrays(1, &shaderProgram.vao)
	gl.DeleteProgram(shaderProgram.program)
}
