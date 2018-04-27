package terrain

import (
	"fmt"
	"math"

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

	min := float64(100)
	max := float64(-100)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			scale := 100.0

			grid[i+j*width] =
				(GetNoise(i, j, scale) +
					GetNoise(i, j, scale/2)/2 +
					GetNoise(i, j, scale/4)/4) / 1.75

			grid[i+j*width] = 1.0 - (math.Pow(grid[i+j*width], 2.8) * 1.5)
			min = math.Min(grid[i+j*width], min)
			max = math.Max(grid[i+j*width], max)
		}
	}

	fmt.Printf("%v, %v\n\n", min, max)
	return grid
}

func GetNoise(x, y int, scale float64) float64 {
	return (terrainGenerator.noise.Eval2(
		float64(x)/scale, float64(y)/scale) + 1) / 2
}
