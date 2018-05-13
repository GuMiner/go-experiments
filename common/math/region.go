package commonMath

type RegionType int

const (
	CircleRegion RegionType = iota
	SquareRegion
)

// Defines a region
type Region struct {
	RegionType RegionType
	Size       float32
}
