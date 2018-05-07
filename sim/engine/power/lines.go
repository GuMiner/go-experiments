package power

import "github.com/go-gl/mathgl/mgl32"

type PowerLine struct {
	start    mgl32.Vec2
	end      mgl32.Vec2
	capacity int64
}
