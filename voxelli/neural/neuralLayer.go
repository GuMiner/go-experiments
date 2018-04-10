package neural

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.
import (
	"math"
	"math/rand"
)

// Defines a deeply-connected neuron with weights to the neurons on the next layer
type neuron struct {
	weights []float32
}

type neuralLayer struct {
	nextLayerSize int
	neurons       []neuron
	bias          float32
}

// Apply the sigmoid function as an activation function
func activationFunction(value float32) float32 {
	if value > 10 {
		return 1.0
	} else if value < -10 {
		return 0.0
	}

	return 1.0 / (1.0 + float32(math.Exp(-float64(value))))
}

// Activates using inputs for these neurons and returns the output for the next layer
// len(inputs) == len(neurons). len(output) == nextLayerSize
func (n *neuralLayer) Evaluate(inputs []float32) []float32 {
	outputs := make([]float32, n.nextLayerSize)
	for i := 0; i < n.nextLayerSize; i++ {
		outputs[i] = n.bias
	}

	for i, input := range inputs {
		for j, weight := range n.neurons[i].weights {
			outputs[j] += weight * input
		}
	}

	for i := 0; i < n.nextLayerSize; i++ {
		outputs[i] = activationFunction(outputs[i])
	}

	return outputs
}

func newNeuralLayer(layerSize, nextLayerSize int) *neuralLayer {
	neurons := make([]neuron, layerSize)
	for i, _ := range neurons {
		neurons[i].weights = make([]float32, nextLayerSize)
	}

	layer := neuralLayer{nextLayerSize: nextLayerSize, neurons: neurons, bias: 1.0}
	layer.Randomize()

	return &layer
}

// Randomizes all weights, excluding the bias.
func (n *neuralLayer) Randomize() {
	for i, _ := range n.neurons {
		for j, _ := range n.neurons[i].weights {
			n.neurons[i].weights[j] = rand.Float32()*2.0 - 1.0
		}
	}
}

func (n *neuralLayer) ProbablyRandomize(randomizeProbability, randomAmount float32) {
	for i, _ := range n.neurons {
		for j, _ := range n.neurons[i].weights {
			if rand.Float32() < randomizeProbability {
				n.neurons[i].weights[j] += rand.Float32()*randomAmount - randomAmount
			}
		}
	}
}

func (n *neuralLayer) CrossMerge(first, second *neuralLayer, crossoverProbability float32) {
	for i, _ := range n.neurons {
		for j, _ := range n.neurons[i].weights {
			if rand.Float32() < crossoverProbability {
				n.neurons[i].weights[j] = second.neurons[i].weights[j]
			} else {
				n.neurons[i].weights[j] = first.neurons[i].weights[j]
			}
		}
	}
}
