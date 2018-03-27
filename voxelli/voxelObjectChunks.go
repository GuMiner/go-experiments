package main

import "fmt"

type MainChunk struct {
}

func (c MainChunk) Name() string {
	return "Main"
}

type PackChunk struct {
	modelCount int
}

func (c PackChunk) Name() string {
	return "Pack"
}

type SizeChunk struct {
	size IntVec3
}

func (c SizeChunk) Name() string {
	return "Size"
}

type VoxelsChunk struct {
	voxels []Voxel
}

func (c VoxelsChunk) Name() string {
	return "Voxels"
}

type PaletteChunk struct {
	palette VoxelPalette
}

func (c PaletteChunk) Name() string {
	return "Palette"
}

type UnknownChunk struct {
	typeName string
}

func (c UnknownChunk) Name() string {
	return fmt.Sprintf("unknown (%v)", c.typeName)
}
