package region

import (
	"go-experiments/common/commonmath"
	"go-experiments/common/commonopengl"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type RegionShadingData struct {
	offset int32
	length int32
}

type RegionShaderProgram struct {
	program uint32

	offsetLoc      int32
	scaleLoc       int32
	orientationLoc int32
	colorLoc       int32

	vao        uint32
	vbo        uint32
	regionData map[commonMath.RegionType]RegionShadingData
}

func (r *RegionShaderProgram) sendRegionDataToShader() {
	offset := int32(0)

	squareRegion := generateRegion(commonMath.SquareRegion)
	r.regionData[commonMath.SquareRegion] = RegionShadingData{
		offset: offset,
		length: int32(len(squareRegion))}

	offset += int32(len(squareRegion))

	triangleRegion := generateRegion(commonMath.TriangleRegion)
	r.regionData[commonMath.TriangleRegion] = RegionShadingData{
		offset: offset,
		length: int32(len(triangleRegion))}

	offset += int32(len(triangleRegion))

	circleRegion := generateRegion(commonMath.CircleRegion)
	r.regionData[commonMath.CircleRegion] = RegionShadingData{
		offset: offset,
		length: int32(len(circleRegion))}

	regions := append(append(squareRegion, triangleRegion...), circleRegion...)

	// 2 -- 2 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(regions)*2*4, gl.Ptr(regions), gl.STATIC_DRAW)
}

func NewRegionShaderProgram() *RegionShaderProgram {
	regionShaderProg := RegionShaderProgram{
		regionData: make(map[commonMath.RegionType]RegionShadingData)}

	// Setup the OpenGL program
	regionShaderProg.program = commonOpenGl.CreateProgram("./ui/shaders/region")

	regionShaderProg.offsetLoc = gl.GetUniformLocation(regionShaderProg.program, gl.Str("offset\x00"))
	regionShaderProg.scaleLoc = gl.GetUniformLocation(regionShaderProg.program, gl.Str("scale\x00"))
	regionShaderProg.orientationLoc = gl.GetUniformLocation(regionShaderProg.program, gl.Str("orientation\x00"))
	regionShaderProg.colorLoc = gl.GetUniformLocation(regionShaderProg.program, gl.Str("givenColor\x00"))

	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &regionShaderProg.vao)
	gl.BindVertexArray(regionShaderProg.vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &regionShaderProg.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, regionShaderProg.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	regionShaderProg.sendRegionDataToShader()

	return &regionShaderProg
}

func (shaderProgram *RegionShaderProgram) PreRender() {
	gl.UseProgram(shaderProgram.program)
	gl.BindVertexArray(shaderProgram.vao)
}

func (shaderProgram *RegionShaderProgram) Render(region *commonMath.Region, color mgl32.Vec3) {
	// Setup location
	gl.Uniform2f(shaderProgram.offsetLoc, region.Position.X(), region.Position.Y())
	gl.Uniform1f(shaderProgram.scaleLoc, region.Scale)
	gl.Uniform1f(shaderProgram.orientationLoc, region.Orientation)
	gl.Uniform3f(shaderProgram.colorLoc, color.X(), color.Y(), color.Z())

	regionShadingData := shaderProgram.regionData[region.RegionType]
	gl.DrawArrays(gl.TRIANGLES, regionShadingData.offset, regionShadingData.length)
}

func (shaderProgram *RegionShaderProgram) Delete() {
	gl.DeleteBuffers(1, &shaderProgram.vbo)
	gl.DeleteVertexArrays(1, &shaderProgram.vao)
	gl.DeleteProgram(shaderProgram.program)
}
