package renderer

import "github.com/go-gl/mathgl/mgl32"

// Defines a renderer that needs camera / projection updates
type Renderer interface {
	UpdateProjection(projection *mgl32.Mat4)
	UpdateCamera(camera *mgl32.Mat4)
}

func UpdateCameras(renderers []Renderer, camera *mgl32.Mat4) {
	for _, renderer := range renderers {
		renderer.UpdateCamera(camera)
	}
}

func UpdateProjections(renderers []Renderer, projection *mgl32.Mat4) {
	for _, renderer := range renderers {
		renderer.UpdateProjection(projection)
	}
}
