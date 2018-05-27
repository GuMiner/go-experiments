package editorEngine

import (
	"fmt"
	"go-experiments/sim/input"
)

// TODO: Modes should really be saved in game state on a per-game basis.
var EngineState State

type KeyToggle struct {
	toggle *bool
	about  string
}

type KeyAction struct {
	action func()
	about  string
}

var globalKeyToggleActions map[input.KeyAssignment]KeyToggle = make(map[input.KeyAssignment]KeyToggle)
var globalKeyActions map[input.KeyAssignment]KeyAction = make(map[input.KeyAssignment]KeyAction)

var addModeKeyActions map[input.KeyAssignment]KeyAction = make(map[input.KeyAssignment]KeyAction)
var addModeKeySubSelectionActions map[input.KeyAssignment]KeyAction = make(map[input.KeyAssignment]KeyAction)
var drawModeKeySubSelectionActions map[input.KeyAssignment]KeyAction = make(map[input.KeyAssignment]KeyAction)

func Init() {
	EngineState = State{
		Mode:             Select,
		InAddMode:        PowerPlant,
		InDrawMode:       TerrainFlatten,
		ItemSubSelection: Item1,
		SnapToGrid:       true,
		SnapToElements:   true,
		SnapToAngle:      false}

	globalKeyToggleActions[input.SnapToGridKey] = KeyToggle{&EngineState.SnapToGrid, "snap to grid"}
	globalKeyToggleActions[input.SnapToAngleKey] = KeyToggle{&EngineState.SnapToAngle, "snap to angle"}
	globalKeyToggleActions[input.SnapToElementsKey] = KeyToggle{&EngineState.SnapToElements, "snap to elements"}

	globalKeyActions[input.SelectModeKey] = KeyAction{func() { EngineState.Mode = Select }, "selection"}
	globalKeyActions[input.AddModeKey] = KeyAction{func() { EngineState.Mode = Add }, "addition"}
	globalKeyActions[input.DrawModeKey] = KeyAction{func() { EngineState.Mode = Draw }, "draw"}

	addModeKeyActions[input.PowerPlantAddModeKey] = KeyAction{func() { EngineState.InAddMode = PowerPlant }, "power plant add"}
	addModeKeyActions[input.PowerLineAddModeKey] = KeyAction{func() { EngineState.InAddMode = PowerLine }, "power line add"}
	addModeKeyActions[input.RoadLineAddModeKey] = KeyAction{func() { EngineState.InAddMode = RoadLine }, "road line add"}

	addModeKeySubSelectionActions[input.ItemAdd1Key] = KeyAction{func() { EngineState.ItemSubSelection = Item1 }, "item 1 selection"}
	addModeKeySubSelectionActions[input.ItemAdd2Key] = KeyAction{func() { EngineState.ItemSubSelection = Item2 }, "item 2 selection"}
	addModeKeySubSelectionActions[input.ItemAdd3Key] = KeyAction{func() { EngineState.ItemSubSelection = Item3 }, "item 3 selection"}
	addModeKeySubSelectionActions[input.ItemAdd4Key] = KeyAction{func() { EngineState.ItemSubSelection = Item4 }, "item 4 selection"}
	addModeKeySubSelectionActions[input.ItemAdd5Key] = KeyAction{func() { EngineState.ItemSubSelection = Item5 }, "item 5 selection"}
	addModeKeySubSelectionActions[input.ItemAdd6Key] = KeyAction{func() { EngineState.ItemSubSelection = Item6 }, "item 6 selection"}

	drawModeKeySubSelectionActions[input.TerrainFlattenKey] = KeyAction{func() { EngineState.InDrawMode = TerrainFlatten }, "terrain flatten"}
	drawModeKeySubSelectionActions[input.TerrainSharpenKey] = KeyAction{func() { EngineState.InDrawMode = TerrainSharpen }, "terrain sharpen"}
	drawModeKeySubSelectionActions[input.TerrainTreesKey] = KeyAction{func() { EngineState.InDrawMode = TerrainTrees }, "terrain trees"}
	drawModeKeySubSelectionActions[input.TerrainShrubsKey] = KeyAction{func() { EngineState.InDrawMode = TerrainShrubs }, "terrain shrubs"}
	drawModeKeySubSelectionActions[input.TerrainHillsKey] = KeyAction{func() { EngineState.InDrawMode = TerrainHills }, "terrain hills"}
	drawModeKeySubSelectionActions[input.TerrainValleysKey] = KeyAction{func() { EngineState.InDrawMode = TerrainValleys }, "terrain valleys"}
}

func checkToggleKeys(toggleKeys map[input.KeyAssignment]KeyToggle) bool {
	updated := false
	for key, keyToggle := range toggleKeys {
		if input.IsTyped(key) {
			*keyToggle.toggle = !*keyToggle.toggle
			fmt.Printf("Toggled '%v' to %v.\n", keyToggle.about, *keyToggle.toggle)
			updated = true
		}
	}

	return updated
}

func checkActionKeys(actionKeys map[input.KeyAssignment]KeyAction) bool {
	updated := false
	for key, keyAction := range actionKeys {
		if input.IsTyped(key) {
			keyAction.action()
			fmt.Printf("Entered '%v' mode.\n", keyAction.about)
			updated = true
		}
	}

	return updated
}

// Update toggles and the current edit mode.
func Update() (updatedSelection, updatedSubSelection bool) {
	updatedSelection = checkToggleKeys(globalKeyToggleActions) || updatedSelection
	updatedSelection = checkActionKeys(globalKeyActions) || updatedSelection
	updatedSubSelection = false

	if EngineState.Mode == Add {
		updatedSelection = checkActionKeys(addModeKeyActions) || updatedSelection
		updatedSubSelection = checkActionKeys(addModeKeySubSelectionActions) || updatedSubSelection
	} else if EngineState.Mode == Draw {
		updatedSubSelection = checkActionKeys(drawModeKeySubSelectionActions) || updatedSubSelection
	}

	return updatedSelection, updatedSubSelection
}
