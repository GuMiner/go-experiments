package neural

// See https://github.com/ArztSamuel/Applying_EANNs for the inspiration for this.
import (
	"math"
	"math/rand"
)

// Defines a deeply-connected neuron with weights to the neurons on the next layer
type Neuron struct {
	weights []float32
}

type NeuralLayer struct {
	nextLayerSize int
	neurons       []Neuron
	bias          float32
}

type NeuralNet struct {
	layers []*NeuralLayer
}

func (n *NeuralNet) Evaluate(inputs []float32) []float32 {
	for _, layer := range n.layers {
		inputs = layer.Evaluate(inputs)
	}

	return inputs
}

func (n *NeuralNet) ProbablyRandomize(randomizeProbability, randomAmount float32) {
	for _, layer := range n.layers {
		layer.ProbablyRandomize(randomizeProbability, randomAmount)
	}
}

func (n *NeuralNet) CrossMerge(first, second *NeuralNet, crossoverProbability float32) {
	for i, layer := range n.layers {
		layer.CrossMerge(first.layers[i], second.layers[i], crossoverProbability)
	}
}

func NewNeuralNet(layerSizes []int, outputCount int) *NeuralNet {
	net := NeuralNet{layers: []*NeuralLayer{}}
	for i, layerSize := range layerSizes {
		nextLayerSize := outputCount
		if i != len(layerSizes)-1 {
			nextLayerSize = layerSizes[i+1]
		}

		net.layers = append(net.layers, newNeuralLayer(layerSize, nextLayerSize))
	}

	return &net
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

func newNeuralLayer(layerSize, nextLayerSize int) *NeuralLayer {
	neurons := make([]Neuron, layerSize)
	for i, _ := range neurons {
		neurons[i].weights = make([]float32, nextLayerSize)
	}

	layer := NeuralLayer{nextLayerSize: nextLayerSize, neurons: neurons, bias: 1.0}
	layer.Randomize()

	return &layer
}

// Randomizes all weights, excluding the bias.
func (n *NeuralLayer) Randomize() {
	for i, _ := range n.neurons {
		for j, _ := range n.neurons[i].weights {
			n.neurons[i].weights[j] = rand.Float32()*2.0 - 1.0
		}
	}
}

func (n *NeuralLayer) ProbablyRandomize(randomizeProbability, randomAmount float32) {
	for i, _ := range n.neurons {
		for j, _ := range n.neurons[i].weights {
			if rand.Float32() < randomizeProbability {
				n.neurons[i].weights[j] += rand.Float32()*randomAmount - randomAmount
			}
		}
	}
}

func (n *NeuralLayer) CrossMerge(first, second *NeuralLayer, crossoverProbability float32) {
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
