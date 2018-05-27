package resource

// Defines in-game resources that are usually always present and
// not configurable as they are part of the core game structure.
type UtilityResource struct {
	Power    float32
	Citizens int // TODO: Citizens have education, age, etc. that affect their work output
	// TODO: Add more well-known resources here.
}
