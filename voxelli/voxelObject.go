package main

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
}

func Load(fileName string) VoxelObject {
	// file := ReadFileAsBytes(fileName)
	// TODO: Copy over and reimplement from the rust variant
	return VoxelObject{}
}
