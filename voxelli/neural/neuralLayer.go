package neural

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.
import (
	"math"
	"math/rand"
)

// Defines a deeply-connected neuron with weights to the neurons on the next layer
type Neuron struct {
	Weights []float32
}

type NeuralLayer struct {
	NextLayerSize int
	Neurons       []Neuron
	Bias          float32
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
func (n *NeuralLayer) Evaluate(inputs []float32) []float32 {
	outputs := make([]float32, n.NextLayerSize)
	for i := 0; i < n.NextLayerSize; i++ {
		outputs[i] = n.Bias
	}

	for i, input := range inputs {
		for j, weight := range n.Neurons[i].Weights {
			outputs[j] += weight * input
		}
	}

	for i := 0; i < n.NextLayerSize; i++ {
		outputs[i] = activationFunction(outputs[i])
	}

	return outputs
}

func newNeuralLayer(layerSize, nextLayerSize int) *NeuralLayer {
	neurons := make([]Neuron, layerSize)
	for i, _ := range neurons {
		neurons[i].Weights = make([]float32, nextLayerSize)
	}

	layer := NeuralLayer{NextLayerSize: nextLayerSize, Neurons: neurons, Bias: 1.0}
	layer.Randomize()

	return &layer
}

// Randomizes all weights, excluding the bias.
func (n *NeuralLayer) Randomize() {
	for i, _ := range n.Neurons {
		for j, _ := range n.Neurons[i].Weights {
			n.Neurons[i].Weights[j] = rand.Float32()*2.0 - 1.0
		}
	}
}

func (n *NeuralLayer) ProbablyRandomize(randomizeProbability, randomAmount float32) {
	for i, _ := range n.Neurons {
		for j, _ := range n.Neurons[i].Weights {
			if rand.Float32() < randomizeProbability {
				n.Neurons[i].Weights[j] += rand.Float32()*randomAmount - randomAmount
			}
		}
	}
}

func (n *NeuralLayer) CrossMerge(first, second *NeuralLayer, crossoverProbability float32) {
	for i, _ := range n.Neurons {
		for j, _ := range n.Neurons[i].Weights {
			if rand.Float32() < crossoverProbability {
				n.Neurons[i].Weights[j] = second.Neurons[i].Weights[j]
			} else {
				n.Neurons[i].Weights[j] = first.Neurons[i].Weights[j]
			}
		}
	}
}
