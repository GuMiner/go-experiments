package engine

import (
	"go-experiments/common/commonmath"

	"github.com/go-gl/mathgl/mgl32"
)

type HypotheticalRegion struct {
	Color  mgl32.Vec3
	Region commonMath.Region
}

type HypotheticalLine struct {
	Color mgl32.Vec3
	Line  [2]mgl32.Vec2
}

// Defines hypothetical regions for drawing and actions
type HypotheticalActions struct {
	Regions []HypotheticalRegion
	Lines   []HypotheticalLine
}

func (h *HypotheticalActions) Reset() {
	h.Regions = []HypotheticalRegion{}
	h.Lines = []HypotheticalLine{}
}

func (h *HypotheticalActions) SetSingleRegion(region HypotheticalRegion) {
	h.Reset()
	h.Regions = []HypotheticalRegion{region}
}
