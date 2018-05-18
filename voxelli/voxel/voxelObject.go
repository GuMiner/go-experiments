package voxel

import (
	"encoding/binary"
	"fmt"
	"math"

	"go-experiments/common/commonio"
	"go-experiments/common/commonmath"
)

// Defines voxel objects
type Voxel struct {
	Position commonMath.IntVec3
	ColorIdx uint8
}

type SubObject struct {
	Voxels []Voxel
}

type VoxelObject struct {
	SubObjects []SubObject
	MinBounds  commonMath.IntVec3
	MaxBounds  commonMath.IntVec3
	Palette    *VoxelPalette
}

// Defines voxel types
type ChunkType interface {
	Name() string
}

func CheckBounds(current *commonMath.IntVec3, comparison commonMath.IntVec3, comparisonFunc func(int, int) int) {
	current[0] = comparisonFunc(current[0], comparison[0])
	current[1] = comparisonFunc(current[1], comparison[1])
	current[2] = comparisonFunc(current[2], comparison[2])
}

func (voxelObject *VoxelObject) ComputeBounds() {
	voxelObject.MinBounds = commonMath.IntVec3{math.MaxInt32, math.MaxInt32, math.MaxInt32}
	voxelObject.MaxBounds = commonMath.IntVec3{math.MinInt32, math.MinInt32, math.MinInt32}

	minFunc := func(x, y int) int {
		if x < y {
			return x
		}

		return y
	}

	maxFunc := func(x, y int) int {
		if x > y {
			return x
		}

		return y
	}

	for _, object := range voxelObject.SubObjects {
		for _, voxel := range object.Voxels {
			CheckBounds(&voxelObject.MinBounds, voxel.Position, minFunc)
			CheckBounds(&voxelObject.MaxBounds, voxel.Position, maxFunc)
		}
	}
}

// Parses out our chunk
func parseChunk(data []uint8) (chunkType ChunkType, bytesRead int) {
	// Validate chunk ID exists
	bytesRead = 4
	checkLength(0, 4, data)

	// Validate the byte sizes exist
	checkLength(4, 8, data)
	bytesRead += 8

	chunkId := string(data[0:4])
	switch chunkId {
	case "MAIN":
		chunkType = MainChunk{}
	case "PACK":
		bytesRead += 4
		checkLength(12, 4, data)

		modelCount := int(binary.LittleEndian.Uint32(data[12:16]))
		chunkType = PackChunk{modelCount: modelCount}
	case "SIZE":
		bytesRead += 12
		checkLength(12, 12, data)

		x := int(binary.LittleEndian.Uint32(data[12:16]))
		y := int(binary.LittleEndian.Uint32(data[16:20]))
		z := int(binary.LittleEndian.Uint32(data[20:24]))
		chunkType = SizeChunk{size: commonMath.IntVec3{x, y, z}}
	case "XYZI":
		bytesRead += 4
		checkLength(12, 4, data)

		voxelCount := int(binary.LittleEndian.Uint32(data[12:16]))
		bytesRead += 4 * voxelCount
		checkLength(16, 4*voxelCount, data)

		voxels := make([]Voxel, voxelCount)
		dataIdx := 16
		for i := 0; i < voxelCount; i++ {
			x := int(data[dataIdx])
			y := int(data[dataIdx+1])
			z := int(data[dataIdx+2])
			c := data[dataIdx+3]
			dataIdx += 4

			voxels[i] = Voxel{Position: commonMath.IntVec3{x, y, z}, ColorIdx: c}
		}

		chunkType = VoxelsChunk{voxels: voxels}
	case "RGBA":
		colorElements := 4 * 256
		bytesRead += colorElements
		checkLength(12, colorElements, data)

		var colors [256]commonMath.Color
		for i := 0; i < len(colors); i++ {
			colors[i][0] = data[12+i*4]
			colors[i][1] = data[12+i*4+1]
			colors[i][2] = data[12+i*4+2]
			colors[i][3] = data[12+i*4+3]
		}

		palette := VoxelPalette{Colors: colors}

		chunkType = PaletteChunk{palette: palette}
	default:
		bytesRead = 12 + int(binary.LittleEndian.Uint32(data[4:8]))
		chunkType = UnknownChunk{typeName: chunkId}
	}

	return
}

func checkLength(start, length int, data []uint8) {
	if start+length > len(data) {
		panic(fmt.Errorf(
			"Not enough data in the voxel input (start: %v, length: %v, data length: %v",
			start, length, len(data)))
	}
}

func checkHeader(data []uint8) {
	header := "VOX "
	checkLength(0, len(header), data)
}

func NewVoxelObject(fileName string) *VoxelObject {
	data := commonIo.ReadFileAsBytes(fileName)

	// Validate this is a VOX file and we start out with our proper MAIN chunk
	checkHeader(data[0:])

	checkLength(4, 4, data[4:])
	fmt.Printf("VOX file %v is version %v\n", fileName, binary.LittleEndian.Uint32(data[4:8]))

	byteIndex := 8
	chunkType, bytesRead := parseChunk(data[byteIndex:])
	fmt.Printf("Read a %v chunk\n", chunkType.Name())

	foundChunk, isMainChunk := chunkType.(MainChunk)
	if !isMainChunk {
		panic("Did not find a main chunk as the first chunk of the file, found " + foundChunk.Name())
	}

	// This is either PACK (1 or more models) or SIZE (the onlyy model.)
	byteIndex += bytesRead
	chunkType, bytesRead = parseChunk(data[byteIndex:])
	fmt.Printf("Read a %v chunk\n", chunkType.Name())

	modelCount := 1
	switch val := chunkType.(type) {
	case PackChunk:
		modelCount = val.modelCount

		// Only advance if this is a PACK chunk so we start parsing our model at the same
		// spot regardless of the presense or lack thereof of PACK
		byteIndex += bytesRead
	}

	var voxelObject VoxelObject
	for i := 0; i < modelCount; i++ {
		chunkType, bytesRead = parseChunk(data[byteIndex:])
		fmt.Printf("Read a %v chunk\n", chunkType.Name())

		switch val := chunkType.(type) {
		case SizeChunk:
			byteIndex += bytesRead
		default:
			panic("Expected the model to start with a SIZE chunk, not " + val.Name())
		}

		chunkType, bytesRead = parseChunk(data[byteIndex:])
		fmt.Printf("Read a %v chunk\n", chunkType.Name())

		switch val := chunkType.(type) {
		case VoxelsChunk:
			voxelObject.SubObjects = append(voxelObject.SubObjects, SubObject{Voxels: val.voxels})
			byteIndex += bytesRead
		default:
			panic("Expected the model to have a XYZI chunk, not " + val.Name())
		}
	}

	voxelObject.Palette = &VoxelPalette{Colors: DefaultPalette}

	// Parse out all the optional chunks, which may give us a palette chunk
	for byteIndex < len(data) {
		chunkType, bytesRead = parseChunk(data[byteIndex:])
		// fmt.Printf("Read a %v chunk\n", chunkType.Name())

		foundChunk, isPaletteChunk := chunkType.(PaletteChunk)
		if isPaletteChunk {
			voxelObject.Palette = &foundChunk.palette
		}

		byteIndex += bytesRead
	}

	voxelObject.ComputeBounds()

	totalVoxels := 0
	for _, subObject := range voxelObject.SubObjects {
		totalVoxels += len(subObject.Voxels)
	}

	fmt.Printf("Voxel Object is bounded by %v and %v with %v total voxels\n", voxelObject.MinBounds, voxelObject.MaxBounds, totalVoxels)

	return &voxelObject
}
