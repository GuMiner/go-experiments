package opengl

// Simplifies OpenGL initialization
import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Capabilities struct {
	MaxTextures    int32
	MaxTextureSize int32
}

var globalCapabilities Capabilities

func logOpenGlInfo() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %v\n\n", version)

	fmt.Printf("Capabilities:\n")

	var intCapability int32
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &intCapability)
	fmt.Printf("  Max Texture Units: %v\n", intCapability)
	globalCapabilities.MaxTextures = intCapability

	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &intCapability)
	fmt.Printf("  Max Texture Size: %v\n", intCapability)
	globalCapabilities.MaxTextureSize = intCapability
}

func GetGlCaps() *Capabilities {
	return &globalCapabilities
}

func InitGlfw() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func ConfigureOpenGl() {
	// Startup OpenGL bindings
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	gl.ClearColor(0.2, 0.5, 1.0, 1.0)

	glfw.SwapInterval(1)

	gl.Disable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.LINE_SMOOTH)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	logOpenGlInfo()
}
