package editorEngine

// Defines the current mode of the in-game editor
type EditorMode int

// Ordering is important! Reordering these changes their numerical values
const (
	Select EditorMode = iota
	Add
)

type State struct {
	Mode EditorMode

	SnapToGrid     bool
	SnapToElements bool
	SnapToAngle    bool
}
