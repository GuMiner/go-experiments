package element

type MapElementType int

// Ordering is important! Reordering these changes their numerical values
const (
	PowerLine MapElementType = iota
	PowerPlant
	// TODO: There's a LOT more to add here. This is enough to get started with putting-and-placing and the like.
)
