package terrain

import (
	"github.com/ojrac/opensimplex-go"
)

type TerrainGenerator struct {
	noise *opensimplex.Noise
}

var terrainGenerator TerrainGenerator

func Init(seed int64) {
	terrainGenerator = TerrainGenerator{
		noise: opensimplex.NewWithSeed(seed)}
}

func Generate(width, height int) []float64 {
	grid := make([]float64, width*height)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			grid[i+j*width] = terrainGenerator.noise.Eval2(
				float64(i), float64(j))
		}
	}

	return grid
}
