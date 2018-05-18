package roadway

import (
	"fmt"
	"go-experiments/common/commonio"
	"go-experiments/common/commonmath"
	"go-experiments/voxelli/config"
	"go-experiments/voxelli/geometry"
	"math"
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

	// Gets the bounds of the road piece.
	GetBounds(gridPos commonMath.IntVec2) []geometry.Intersectable
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

func getGridIdx(position mgl32.Vec2) commonMath.IntVec2 {
	return commonMath.IntVec2{
		int(position.X() / float32(GetGridSize())),
		int(position.Y() / float32(GetGridSize()))}
}

func getGridRelativePos(gridIdx commonMath.IntVec2, position mgl32.Vec2) mgl32.Vec2 {
	return position.Sub(mgl32.Vec2{float32(gridIdx.X() * GetGridSize()), float32(gridIdx.Y() * GetGridSize())})
}

func getRealPosition(position mgl32.Vec2) mgl32.Vec2 {
	return position.Sub(mgl32.Vec2{float32(GetGridSize() / 2), float32(GetGridSize() / 2)})
}

// Defines a 2D roadway with road elements
type Roadway struct {
	roadElements [][]Road
	roadParts    []geometry.Intersectable
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

func (r *Roadway) GetRoadElementIdx(position mgl32.Vec2) commonMath.IntVec2 {
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

// Given positions and directions, gets the distance and normal of each intersection.
func (r *Roadway) GetBoundaries(positions []mgl32.Vec2, directions []mgl32.Vec2) ([]float32, []mgl32.Vec2) {
	boundaryLengths := make([]float32, len(positions))
	normals := make([]mgl32.Vec2, len(positions))

	for i, position := range positions {
		if !r.InBounds(position) {
			boundaryLengths[i] = 0.0
		}

		offsetPosition := getOffsetPosition(position)
		vector := geometry.NewVector(offsetPosition, directions[i])

		// Brute-force find the closest intersection by looking at the *entire* roadway
		hasIntersection := false
		minIntersectionDist := float32(math.MaxFloat32)
		minNormal := mgl32.Vec2{0, 0}
		for _, intersectable := range r.roadParts {
			if intersects, intersectionPoint, intersectionNormal := intersectable.Intersects(vector); intersects {
				realPos := getRealPosition(intersectionPoint)
				intersectionLen := realPos.Sub(position).Len()
				if intersectionLen < minIntersectionDist {
					minIntersectionDist = intersectionLen
					minNormal = intersectionNormal
					hasIntersection = true
				}
			}
		}

		if hasIntersection {
			boundaryLengths[i] = minIntersectionDist
			normals[i] = minNormal
		} else {
			boundaryLengths[i] = 10000.0 // External bounding box
			normals[i] = mgl32.Vec2{1, 0}
		}
	}

	return boundaryLengths, normals
}

// Defines the 2D bounds of road elements
func GetGridSize() int {
	return config.Config.Simulation.RoadwaySize
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

	file := commonIo.ReadFile(fileName)

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

			newRoad := newRoad(roadType, optionalData)
			roadway.roadElements[xSize-(i+1)][j] = newRoad
			roadway.roadParts = append(roadway.roadParts, newRoad.GetBounds(commonMath.IntVec2{xSize - (i + 1), j})...)
		}
	}

	fmt.Printf("Roadway size: [%v, %v]\n\n", len(roadway.roadElements), len(roadway.roadElements[0]))
	return &roadway
}
