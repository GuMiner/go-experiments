package main

import (
	"encoding/binary"
	"fmt"
)

// Defines voxel objects
type Voxel struct {
	position IntVec3
	colorIdx uint8
}

type SubObject struct {
	voxels []Voxel
}

type VoxelObject struct {
	subObjects []SubObject
	minBounds  IntVec3
	maxBounds  IntVec3
	palette    *VoxelPalette
}

// Defines voxel types
type ChunkType interface {
	Name() string
	Add(voxelObject *VoxelObject)
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
		chunkType = SizeChunk{size: IntVec3{x, y, z}}
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
	data := ReadFileAsBytes(fileName)

	checkHeader(data[0:])

	checkLength(4, 4, data[4:])
	fmt.Printf("VOX file %v is version %v\n", fileName, binary.LittleEndian.Uint32(data[4:8]))

	byteIndex := 8
	chunkType, bytesRead := parseChunk(data[byteIndex:])
	fmt.Printf("Read a %v chunk\n", chunkType.Name())

	_, isMainChunk := chunkType.(MainChunk)
	if !isMainChunk {
		panic("Did not find a main chunk as the first chunk of the file!")
	}

	var voxelObject VoxelObject

	byteIndex += bytesRead
	chunkType, bytesRead = parseChunk(data[byteIndex:])
	fmt.Printf("Read a %v chunk\n", chunkType.Name())

	switch val := chunkType.(type) {
	case MainChunk:
		panic("Did not expect to read in more than one main chunk!")
	case PackChunk:
		val.Add(&voxelObject)
	}

	return &voxelObject
}
