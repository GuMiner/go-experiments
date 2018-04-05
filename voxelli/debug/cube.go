package debug

// Defines a small cube
import (
	"go-experiments/voxelli/opengl"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Directly from https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
var cubeVertices = []mgl32.Vec3{
	// Bottom
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},

	// Top
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},

	// Front
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},

	// Back
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},

	// Left
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, -0.5},
	mgl32.Vec3{-0.5, -0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, 0.5},
	mgl32.Vec3{-0.5, 0.5, -0.5},

	// Right
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, -0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, -0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, -0.5},
	mgl32.Vec3{0.5, 0.5, 0.5}}

var cubeColorVertices = []mgl32.Vec3{
	// Bottom
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.5, 0.2},

	// Top
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.2},

	// Front
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.2},

	// Back
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.5},

	// Left
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.5},
	mgl32.Vec3{0.5, 0.5, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.2},
	mgl32.Vec3{0.5, 0.2, 0.5},

	// Right
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.5, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.5, 0.2},
	mgl32.Vec3{0.2, 0.2, 0.5},
	mgl32.Vec3{0.2, 0.2, 0.2}}

type Cube struct {
	shaderProgram uint32

	projectionLoc    int32
	cameraLoc        int32
	modelLoc         int32
	timeLoc          int32
	colorOverrideLoc int32

	vao         uint32
	positionVbo uint32
	colorVbo    uint32
}

var cube *Cube

func InitCube() {
	cube = new(Cube)
	cube.shaderProgram = opengl.CreateProgram("./shaders/basicRenderer")

	// Get locations of everything used in this program.
	cube.projectionLoc = gl.GetUniformLocation(cube.shaderProgram, gl.Str("projection\x00"))
	cube.cameraLoc = gl.GetUniformLocation(cube.shaderProgram, gl.Str("camera\x00"))
	cube.modelLoc = gl.GetUniformLocation(cube.shaderProgram, gl.Str("model\x00"))
	cube.timeLoc = gl.GetUniformLocation(cube.shaderProgram, gl.Str("runTime\x00"))
	cube.colorOverrideLoc = gl.GetUniformLocation(cube.shaderProgram, gl.Str("colorOverride\x00"))

	// Setup triangles for us to draw
	gl.GenVertexArrays(1, &cube.vao)
	gl.BindVertexArray(cube.vao)

	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &cube.positionVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// 3 -- 3 floats / vertex. 4 -- float32
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*3*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(1)
	gl.GenBuffers(1, &cube.colorVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.colorVbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.BufferData(gl.ARRAY_BUFFER, len(cubeColorVertices)*3*4, gl.Ptr(cubeColorVertices), gl.STATIC_DRAW)
}

func vertexCount() int32 {
	return int32(len(cubeVertices))
}

func Render(time float32, color mgl32.Vec4, model *mgl32.Mat4) {
	gl.UseProgram(cube.shaderProgram)

	gl.Uniform1f(cube.timeLoc, time)
	gl.Uniform4fv(cube.colorOverrideLoc, 1, &color[0])
	gl.UniformMatrix4fv(cube.modelLoc, 1, false, &model[0])

	gl.BindVertexArray(cube.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount())
}

func GetCube() *Cube {
	return cube
}

func (renderer *Cube) UpdateProjection(projection *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.projectionLoc, 1, false, &projection[0])
}

func (renderer *Cube) UpdateCamera(camera *mgl32.Mat4) {
	gl.UseProgram(renderer.shaderProgram)
	gl.UniformMatrix4fv(renderer.cameraLoc, 1, false, &camera[0])
}

func DeleteCube() {
	gl.DeleteBuffers(1, &cube.positionVbo)
	gl.DeleteBuffers(1, &cube.colorVbo)
	gl.DeleteVertexArrays(1, &cube.vao)
	gl.DeleteProgram(cube.shaderProgram)
	cube = nil
}
