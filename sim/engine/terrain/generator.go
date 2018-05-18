package terrain

import (
	"fmt"
	"go-experiments/common/commonmath"
	"go-experiments/sim/config"
	"math"

	"github.com/ojrac/opensimplex-go"
)

type TerrainGenerator struct {
	noise *opensimplex.Noise

	// Offset values used to place the noise in the 0 to 1 range.
	hasSetOffsetFactors bool
	linearOffset        float32
	scaleOffset         float32
}

var terrainGenerator TerrainGenerator

func Init(seed int) {
	terrainGenerator = TerrainGenerator{
		noise:               opensimplex.NewWithSeed(int64(seed)),
		hasSetOffsetFactors: false}
}

func Generate(width, height, xOffset, yOffset int) []float32 {
	grid := make([]float32, width*height)

	min := float32(1e10)
	max := float32(-1e10)

	generation := config.Config.Terrain.Generation

	// Generate random noise values
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			grid[i+j*width] = getNoise(i+xOffset, j+yOffset, generation.MaxNoiseScale)*generation.MaxNoiseContribution +
				getNoise(i+xOffset, j+yOffset, generation.MedNoiseScale)*generation.MedNoiseContribution +
				getNoise(i+xOffset, j+yOffset, generation.MinNoiseScale)*generation.MinNoiseContribution

			min = commonMath.MinFloat32(grid[i+j*width], min)
			max = commonMath.MaxFloat32(grid[i+j*width], max)
		}
	}

	// We only set offset factors once, so that edges of other random grids align.
	if !terrainGenerator.hasSetOffsetFactors {
		terrainGenerator.linearOffset = -min
		terrainGenerator.scaleOffset = 1.0 / (max - min)

		terrainGenerator.hasSetOffsetFactors = true
	}

	// Rescale and apply the power factor to flatten lowlands
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			grid[i+j*width] = (grid[i+j*width] + terrainGenerator.linearOffset) * terrainGenerator.scaleOffset
			grid[i+j*width] = float32(math.Pow(float64(grid[i+j*width]), float64(generation.PowerFactor)))
		}
	}

	fmt.Printf("Generated %v random elements within [%v, %v]\n", width*height, min, max)
	return grid
}

func getNoise(x, y int, scale float32) float32 {
	return float32(terrainGenerator.noise.Eval2(float64(x)/float64(scale), float64(y)/float64(scale)))
}
