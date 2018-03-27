package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

// Defines the road types that exist in our file.
type RoadType int

const (
	NotARoadType RoadType = iota
	StraightRoadType
	CurvedRoadType
	MaxRoadType
)

// Defines a generic road element
type Road interface {
	// Returns true if the position is in-bounds on the road.
	// position: Normalized from (0, 0) to GetRoadBounds(), guaranteed to be within the road piece
	InBounds(position mgl32.Vec2) bool
}

// Defines a 2D roadway with road elements
type Roadway struct {
	roadElements [][]Road
}

func NewRoad(roadType RoadType, optionalData int) Road {
	switch roadType {
	case StraightRoadType:
		return StraightRoad{rotated: optionalData != 0}
	case CurvedRoadType:
		return CurvedRoad{rotation: optionalData}
	case NotARoadType:
		return OutOfBoundsRoad{}
	default:
		return OutOfBoundsRoad{}
	}
}

// Defines the 2D bounds of road elements
func GetGridSize() int {
	return 40 // All models are 40x40
}

func ParseInt(item string) int {
	i, err := strconv.ParseInt(item, 10, 32)
	if err != nil {
		panic(err)
	}

	return int(i)
}

func ParseRoadType(item string) RoadType {
	roadType := ParseInt(item)
	if roadType < 0 || roadType >= int(MaxRoadType) {
		panic(fmt.Sprintf("Did not parse a road type, parsed %v instead.", roadType))
	}

	return RoadType(roadType)
}

func NewRoadway(fileName string) *Roadway {
	newlineSplitFunction := func(c rune) bool {
		return c == '\n' || c == '\r'
	}

	spaceSplitFunction := func(c rune) bool {
		return c == ' ' || c == '\t'
	}

	file := ReadFile(fileName)

	lines := strings.FieldsFunc(file, newlineSplitFunction)
	if len(lines) < 3 {
		panic("Expected at least three lines in the file, not enough found.")
	}

	// The roadway format corresponds to what you see when you edit it.
	// *empty* newlines are allowed anywhere in the format
	// Line 1: xSize
	// Line 2: ySize
	xSize := ParseInt(lines[0])
	ySize := ParseInt(lines[1])
	fmt.Printf("Started parsing a roadway grid of size [%v, %v]\n", xSize, ySize)

	// Line 3-ySize + 3: Y lines, flipped upside down to match the screen display
	// Any number of spaces or tabs can be used for item delimiters
	// Items are defined as RoadType:OptionalValue
	if len(lines) < 3+ySize {
		panic(fmt.Sprintf("Did not find enough lines to parse the full roadway grid, Found %v, expected %v", len(lines), ySize+3))
	}

	roadway := Roadway{roadElements: make([][]Road, xSize)}
	for i, _ := range roadway.roadElements {
		roadway.roadElements[i] = make([]Road, ySize)
	}

	for i, line := range lines[2:] {
		roadParts := strings.FieldsFunc(line, spaceSplitFunction)
		if len(roadParts) != xSize {
			panic(fmt.Sprintf("Found %v road elements, expected %v", len(roadParts), xSize))
		}

		for j, part := range roadParts {
			subParts := strings.Split(part, ":")
			roadType := ParseRoadType(subParts[0])

			optionalData := 0 // Default if the item takes in optional data but does not use it.
			if len(subParts) > 1 {
				// We have optional data. Create it
				optionalData = ParseInt(subParts[1])
			}

			roadway.roadElements[i][ySize-(j+1)] = NewRoad(roadType, optionalData)
		}
	}

	return &roadway
}
