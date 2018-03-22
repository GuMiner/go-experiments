package main

import "fmt"

type MainChunk struct {
}

func (c MainChunk) Name() string {
	return "Main"
}
func (c MainChunk) Add(obj *VoxelObject) {
}

type PackChunk struct {
	modelCount int
}

func (c PackChunk) Name() string {
	return "Pack"
}
func (c PackChunk) Add(obj *VoxelObject) {
}

type SizeChunk struct {
	size IntVec3
}

func (c SizeChunk) Name() string {
	return "Size"
}
func (c SizeChunk) Add(obj *VoxelObject) {
}

type VoxelsChunk struct {
	voxels []Voxel
}

func (c VoxelsChunk) Name() string {
	return "Voxels"
}
func (chunk VoxelsChunk) Add(obj *VoxelObject) {
	obj.subObjects = append(obj.subObjects, SubObject{voxels: chunk.voxels})
}

type PaletteChunk struct {
	palette VoxelPalette
}

func (c PaletteChunk) Name() string {
	return "Palette"
}
func (chunk PaletteChunk) Add(obj *VoxelObject) {
	*obj.palette = chunk.palette
}

type UnknownChunk struct {
	typeName string
}

func (c UnknownChunk) Name() string {
	return fmt.Sprintf("unknown (%v)", c.typeName)
}
func (c UnknownChunk) Add(obj *VoxelObject) {
}
