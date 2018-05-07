package editorEngine

// Defines the current mode of the in-game editor
type EdtiorMode int

// Ordering is important! Reordering these changes their numerical values
const (
	Select EdtiorMode = iota
	Add
)

type State struct {
	Mode EdtiorMode

	SnapToGrid     bool
	SnapToElements bool
	SnapToAngle    bool
}
