package editorEngine

// Defines the current mode of the in-game editor
type EditorMode int

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

// Defines the mode of the power plant add action
type EditorPlantAddMode int

const (
	CoalPlant EditorPlantAddMode = iota
	NuclearPlant
	NaturalGasPlant
	WindPlant
	SolarPlant
	GeothermalPlant
)

type State struct {
	Mode                EditorMode
	InAddMode           EditorAddMode
	InPowerPlantAddMode EditorPlantAddMode

	SnapToGrid     bool
	SnapToElements bool
	SnapToAngle    bool
}
