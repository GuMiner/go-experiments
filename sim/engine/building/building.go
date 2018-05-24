package building

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Building struct {
	location    mgl32.Vec2
	size        float32
	orientation float32

	buildingType    string
	storedResources map[string]float32

	gridId int
}
