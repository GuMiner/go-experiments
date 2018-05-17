package editorEngine

type EditorMode int

const (
	Select EditorMode = iota
	Add
	Draw
)

type EditorAddMode int

const (
	PowerPlant EditorAddMode = iota
	PowerLine
)

type EditorPlantAddMode int

const (
	CoalPlant EditorPlantAddMode = iota
	NuclearPlant
	NaturalGasPlant
	WindPlant
	SolarPlant
	GeothermalPlant
)

type EditorDrawMode int

const (
	TerrainFlatten EditorDrawMode = iota
	TerrainSharpen
	TerrainTrees
	TerrainShrubs
	TerrainHills
	TerrainValleys
)

type State struct {
	Mode                EditorMode
	InAddMode           EditorAddMode
	InDrawMode          EditorDrawMode
	InPowerPlantAddMode EditorPlantAddMode

	SnapToGrid     bool
	SnapToElements bool
	SnapToAngle    bool
}
