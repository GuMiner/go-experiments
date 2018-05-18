package commonOpenGl

// Simplifies OpenGL initialization
import (
	"fmt"
	"go-experiments/common/commonconfig"
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
	glfw.WindowHint(glfw.ContextVersionMajor, commonConfig.Config.Window.OpenGlMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, commonConfig.Config.Window.OpenGlMinor)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if commonConfig.Config.Window.Samples > 0 {
		glfw.WindowHint(glfw.Samples, commonConfig.Config.Window.Samples)
	}
}

func ConfigureOpenGl() {
	// Startup OpenGL bindings
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	gl.ClearColor(
		commonConfig.Config.Window.BackgroundColor.X,
		commonConfig.Config.Window.BackgroundColor.Y,
		commonConfig.Config.Window.BackgroundColor.Z,
		1.0)

	glfw.SwapInterval(1)

	if commonConfig.Config.Window.Samples > 0 {
		gl.Enable(gl.MULTISAMPLE)
	} else {
		gl.Disable(gl.MULTISAMPLE)
	}

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.LINE_SMOOTH)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	logOpenGlInfo()
}
