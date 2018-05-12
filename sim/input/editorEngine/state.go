package editorEngine

// Defines the current mode of the in-game editor
type EditorMode int

// Ordering is important! Reordering these changes their numerical values
const (
	Select EditorMode = iota
	Add
)

// Defines the mode of the add action
type EditorAddMode int

const (
	PowerPlant EditorAddMode = iota
	PowerLine
)

type State struct {
	Mode      EditorMode
	InAddMode EditorAddMode

	SnapToGrid     bool
	SnapToElements bool
	SnapToAngle    bool
}
