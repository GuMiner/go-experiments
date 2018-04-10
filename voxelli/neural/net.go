package neural

type NeuralNet struct {
	layers []*neuralLayer
}

// Evaluates the neural net, returning the output. The input count must be the first layer size.
func (n *NeuralNet) Evaluate(inputs []float32) []float32 {
	for _, layer := range n.layers {
		inputs = layer.Evaluate(inputs)
	}

	return inputs
}

// Randomizes all neurons in the net by the random amount with the given randomize probability
func (n *NeuralNet) ProbablyRandomize(randomizeProbability, randomAmount float32) {
	for _, layer := range n.layers {
		layer.ProbablyRandomize(randomizeProbability, randomAmount)
	}
}

// Merges two neural nets into the current neural network, using the given probability to grab items from the second neural net instead of the first.
func (n *NeuralNet) CrossMerge(first, second *NeuralNet, crossoverProbability float32) {
	for i, layer := range n.layers {
		layer.CrossMerge(first.layers[i], second.layers[i], crossoverProbability)
	}
}

// Generates a fully-connected neural net, with inputs == layerSizes[0], outputs == outputCount.
// Neurons are 32-bit floating point values with a constant (1.0) bias applied at each layer.
func NewNeuralNet(layerSizes []int, outputCount int) *NeuralNet {
	net := NeuralNet{layers: []*neuralLayer{}}
	for i, layerSize := range layerSizes {
		nextLayerSize := outputCount
		if i != len(layerSizes)-1 {
			nextLayerSize = layerSizes[i+1]
		}

		net.layers = append(net.layers, newNeuralLayer(layerSize, nextLayerSize))
	}

	return &net
}
