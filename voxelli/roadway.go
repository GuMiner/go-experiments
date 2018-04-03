package main

import (
	"fmt"
	"go-experiments/voxelli/utils"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

// Defines the road types that exist in our file.
type RoadType int

// Ordering is important! Road types are parsed according to these constant values.
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

	// Given a position and direction on the piece, finds the piece boundary.
	// If the boundary is out-of-bounds, returns (true, {0, 0}, and the relative position of the intersection)
	// If the boundary leads to another piece, returns (false, the offset to the next grid pos, and the relative position of the intersection on that next piece)
	FindBoundary(position mgl32.Vec2, direction mgl32.Vec2) (bool, utils.IntVec2, mgl32.Vec2)
}

func newRoad(roadType RoadType, optionalData int) Road {
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

func getOffsetPosition(position mgl32.Vec2) mgl32.Vec2 {
	return position.Add(mgl32.Vec2{float32(GetGridSize() / 2), float32(GetGridSize() / 2)})
}

func getGridIdx(position mgl32.Vec2) utils.IntVec2 {
	return utils.IntVec2{
		int(position.X() / float32(GetGridSize())),
		int(position.Y() / float32(GetGridSize()))}
}

func getGridRelativePos(gridIdx utils.IntVec2, position mgl32.Vec2) mgl32.Vec2 {
	return position.Sub(mgl32.Vec2{float32(gridIdx.X() * GetGridSize()), float32(gridIdx.Y() * GetGridSize())})
}

func getRealPosition(gridIdx utils.IntVec2, gridRelativePos mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		float32(gridIdx.X()*GetGridSize()) + gridRelativePos.X() - float32(GetGridSize()/2),
		float32(gridIdx.Y()*GetGridSize()) + gridRelativePos.Y() - float32(GetGridSize()/2)}
}

// Defines a 2D roadway with road elements
type Roadway struct {
	roadElements [][]Road
}

func (r *Roadway) InBounds(position mgl32.Vec2) bool {
	offsetPos := getOffsetPosition(position)

	// Outside of the entire road grid
	if offsetPos.X() < 0 || offsetPos.Y() < 0 || offsetPos.X() > float32(GetGridSize()*len(r.roadElements)) || offsetPos.Y() > float32(GetGridSize()*len(r.roadElements[0])) {
		return false
	}

	gridIdx := getGridIdx(offsetPos)
	offsetPos = getGridRelativePos(gridIdx, offsetPos)

	return r.roadElements[gridIdx.X()][gridIdx.Y()].InBounds(offsetPos)
}

func (r *Roadway) GetRoadElementIdx(position mgl32.Vec2) utils.IntVec2 {
	offsetPos := getOffsetPosition(position)
	return getGridIdx(offsetPos)
}

func (r *Roadway) InAllBounds(positions []mgl32.Vec2) bool {
	for _, position := range positions {
		if !r.InBounds(position) {
			return false
		}
	}

	return true
}

func (r *Roadway) GetBoundaries(positions []mgl32.Vec2, directions []mgl32.Vec2) []float32 {
	boundaryLengths := make([]float32, len(positions))

	for i, position := range positions {
		if !r.InBounds(position) {
			boundaryLengths[i] = 0.0
		}

		offsetPosition := getOffsetPosition(position)
		gridIdx := getGridIdx(offsetPosition)
		gridRelativePos := getGridRelativePos(gridIdx, offsetPosition)

		// Iterate through roadway pieces until we find an intersection
		isBoundary, gridIdxOffset, intersectionPos := r.roadElements[gridIdx.X()][gridIdx.Y()].FindBoundary(gridRelativePos, directions[i])
		for !isBoundary {
			// TODO: I probably need to mess with this a bit.
			gridIdx := utils.IntVec2{gridIdx.X() + gridIdxOffset.X(), gridIdx.Y() + gridIdxOffset.Y()}
			isBoundary, gridIdxOffset, intersectionPos = r.roadElements[gridIdx.X()][gridIdx.Y()].FindBoundary(intersectionPos, directions[i])
		}

		// Convert that intersection back to a real position and return it.
		realPos := getRealPosition(gridIdx, intersectionPos)
		boundaryLengths[i] = realPos.Sub(position).Len()
	}

	return boundaryLengths
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

	file := utils.ReadFile(fileName)

	lines := strings.FieldsFunc(file, newlineSplitFunction)
	if len(lines) < 3 {
		panic("Expected at least three lines in the file, not enough found.")
	}

	fmt.Printf("Roadway information: %v\n", lines[0])

	// The roadway format corresponds to what you see when you edit it.
	// The first line can be a comment. No other comments are allowed.
	// *empty* newlines are allowed anywhere in the format
	// Line 1: xSize
	// Line 2: ySize
	xSize := ParseInt(lines[1])
	ySize := ParseInt(lines[2])
	fmt.Printf("Started parsing a roadway grid of size [%v, %v]\n", xSize, ySize)

	// Remaining lines: Y lines, flipped upside down to match the screen display
	// Any number of spaces or tabs can be used for item delimiters
	// Items are defined as RoadType:OptionalValue
	if len(lines) < 3+ySize {
		panic(fmt.Sprintf("Did not find enough lines to parse the full roadway grid, Found %v, expected %v", len(lines), ySize+3))
	}

	roadway := Roadway{roadElements: make([][]Road, xSize)}
	for i, _ := range roadway.roadElements {
		roadway.roadElements[i] = make([]Road, ySize)
	}

	for j, line := range lines[3:] {
		roadParts := strings.FieldsFunc(line, spaceSplitFunction)
		if len(roadParts) != xSize {
			panic(fmt.Sprintf("Found %v road elements, expected %v", len(roadParts), xSize))
		}

		for i, part := range roadParts {
			subParts := strings.Split(part, ":")
			roadType := ParseRoadType(subParts[0])

			optionalData := 0 // Default if the item takes in optional data but does not use it.
			if len(subParts) > 1 {
				// We have optional data. Create it
				optionalData = ParseInt(subParts[1])
			}

			roadway.roadElements[xSize-(i+1)][j] = newRoad(roadType, optionalData)
		}
	}

	return &roadway
}
