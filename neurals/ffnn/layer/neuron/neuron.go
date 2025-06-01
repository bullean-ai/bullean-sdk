package neuron

import (
	"github.com/bullean-ai/bullean-sdk/neurals/domain"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/layer/neuron/synapse"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/layer/neuron/synapse/activation"
)

// Neuron is a neurals network node
type Neuron struct {
	A     domain.ActivationType
	In    []*synapse.Synapse
	Out   []*synapse.Synapse
	Value float64
	Index int
}

// NewNeuron returns a neuron with the given activation
func NewNeuron(activation domain.ActivationType) *Neuron {
	return &Neuron{
		A: activation,
	}
}

func (n *Neuron) Fire() {
	var sum float64
	for _, s := range n.In {
		sum += s.Out
	}
	n.Value = n.Activate(sum)

	nVal := n.Value
	for _, s := range n.Out {
		s.Fire(nVal)
	}
}

// Activate applies the neurons activation
func (n *Neuron) Activate(x float64) float64 {
	return activation.GetActivation(n.A).F(x)
}

// DActivate applies the derivative of the neurons activation
func (n *Neuron) DActivate(x float64) float64 {
	return activation.GetActivation(n.A).Df(x)
}
