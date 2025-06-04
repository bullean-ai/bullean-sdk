package domain

import (
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer/neuron/synapse"
	"math/rand"
)

var DefaultFFNNConfig = func(input_len int) *Config {
	return &Config{
		Inputs:     input_len,
		Layout:     []int{30, 70, 70, 70, 70, 70, 70, 30, 2},
		Activation: ActivationSigmoid,
		Mode:       ModeMultiClass,
		Weight:     synapse.NewNormal(1e-20, 1e-20),
		Bias:       true,
	}
}

type Neural struct {
	Model   IModel
	Trainer ITrainer
}

// Differentiable is an activation function and its first order derivative,
// where the latter is expressed as a function of the former for efficiency
type Differentiable interface {
	F(float64) float64
	Df(float64) float64
}

// Mode denotes inference mode
type Mode int

const (
	// ModeDefault is unspecified mode
	ModeDefault Mode = 0
	// ModeMultiClass is for one-hot encoded classification, applies softmax output layer
	ModeMultiClass Mode = 1
	// ModeRegression is regression, applies linear output layer
	ModeRegression Mode = 2
	// ModeBinary is binary classification, applies sigmoid output layer
	ModeBinary Mode = 3
	// ModeMultiLabel is for multilabel classification, applies sigmoid output layer
	ModeMultiLabel Mode = 4
)

// ActivationType is represents a neuron activation function
type ActivationType int

const (
	// ActivationNone is no activation
	ActivationNone ActivationType = 0
	// ActivationSigmoid is a sigmoid activation
	ActivationSigmoid ActivationType = 1
	// ActivationTanh is hyperbolic activation
	ActivationTanh ActivationType = 2
	// ActivationReLU is rectified linear unit activation
	ActivationReLU ActivationType = 3
	// ActivationLinear is linear activation
	ActivationLinear ActivationType = 4
	// ActivationSoftmax is a softmax activation (per layer)
	ActivationSoftmax ActivationType = 5
)

// Config defines the network topology, activations, losses etc
type Config struct {
	// Number of inputs
	Inputs int
	// Defines topology:
	// For instance, [5 3 3] signifies a network with two hidden layers
	// containing 5 and 3 nodes respectively, followed an output layer
	// containing 3 nodes.
	Layout []int
	// Activation functions: {ActivationTanh, ActivationReLU, ActivationSigmoid}
	Activation ActivationType
	// Solver modes: {ModeRegression, ModeBinary, ModeMultiClass, ModeMultiLabel}
	Mode Mode
	// Initializer for weights: {NewNormal(σ, μ), NewUniform(σ, μ)}
	Weight synapse.WeightInitializer `json:"-"`
	// Loss functions: {LossCrossEntropy, LossBinaryCrossEntropy, LossMeanSquared}
	Loss LossType
	// Apply bias nodes
	Bias bool
}

// LossType represents a loss function
type LossType int

func (l LossType) String() string {
	switch l {
	case LossCrossEntropy:
		return "CE"
	case LossBinaryCrossEntropy:
		return "BinCE"
	case LossMeanSquared:
		return "MSE"
	}
	return "N/A"
}

const (
	// LossNone signifies unspecified loss
	LossNone LossType = 0
	// LossCrossEntropy is cross entropy loss
	LossCrossEntropy LossType = 1
	// LossBinaryCrossEntropy is the special case of binary cross entropy loss
	LossBinaryCrossEntropy LossType = 2
	// LossMeanSquared is MSE
	LossMeanSquared LossType = 3
)

// Example is an input-target pair
type Example struct {
	Input    []float64
	Response []float64
}

// Examples is a set of input-output pairs
type Examples []Example

// Shuffle shuffles slice in-place
func (e Examples) Shuffle() {
	for i := range e {
		j := rand.Intn(i + 1)
		e[i], e[j] = e[j], e[i]
	}
}

// Split assigns each element to two new slices
// according to probability p
func (e Examples) Split(p float64) (first, second Examples) {
	for i := 0; i < len(e); i++ {
		if p > rand.Float64() {
			first = append(first, e[i])
		} else {
			second = append(second, e[i])
		}
	}
	return
}

// SplitSize splits slice into parts of size
func (e Examples) SplitSize(size int) []Examples {
	res := make([]Examples, 0)
	for i := 0; i < len(e); i += size {
		res = append(res, e[i:min(i+size, len(e))])
	}
	return res
}

// SplitN splits slice into n parts
func (e Examples) SplitN(n int) []Examples {
	res := make([]Examples, n)
	for i, el := range e {
		res[i%n] = append(res[i%n], el)
	}
	return res
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
