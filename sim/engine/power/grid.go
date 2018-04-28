package power

type PowerGrid struct {
	grid *ResizeableGraph
}

func NewPowerGrid() *PowerGrid {
	grid := PowerGrid{
		grid: NewResizeableGraph()}

	return &grid
}
