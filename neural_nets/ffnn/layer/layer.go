package layer

import (
	"fmt"
	"github.com/bullean-ai/bullean-go/neurals/domain"
	"github.com/bullean-ai/bullean-go/neurals/ffnn/layer/neuron"
	"github.com/bullean-ai/bullean-go/neurals/ffnn/layer/neuron/synapse"
	"github.com/bullean-ai/bullean-go/neurals/ffnn/layer/neuron/synapse/activation"
)

// Layer is a set of neurons and corresponding activation
type Layer struct {
	Neurons []*neuron.Neuron
	A       domain.ActivationType
}

// NewLayer creates a new layer with n nodes
func NewLayer(n int, activation domain.ActivationType) *Layer {
	neurons := make([]*neuron.Neuron, n)

	for i := 0; i < n; i++ {
		act := activation
		neurons[i] = neuron.NewNeuron(act)
	}
	return &Layer{
		Neurons: neurons,
		A:       activation,
	}
}

func (l *Layer) Fire() {
	for _, n := range l.Neurons {
		n.Fire()
	}
	if l.A == domain.ActivationSoftmax {
		outs := make([]float64, len(l.Neurons))
		for i, neuron := range l.Neurons {
			outs[i] = neuron.Value
		}
		sm := activation.Softmax(outs)
		for i, neuron := range l.Neurons {
			neuron.Value = sm[i]
		}
	}
}

// Connect fully connects layer l to next, and initializes each
// synapse with the given weight function
func (l *Layer) Connect(next *Layer, weight synapse.WeightInitializer) {
	synapseId := 0
	for i := range l.Neurons {
		for j := range next.Neurons {
			syn := synapse.NewSynapse(weight())
			syn.Id = synapseId
			l.Neurons[i].Out = append(l.Neurons[i].Out, syn)
			next.Neurons[j].In = append(next.Neurons[j].In, syn)
			synapseId += 1
		}
	}
}

func (l *Layer) ConnectPrepared(next *Layer) {
	for i := range l.Neurons {
		for j := range next.Neurons {
			for k := 0; k < len(l.Neurons[i].Out); k++ {
				for m := 0; m < len(next.Neurons[j].In); m++ {
					if l.Neurons[i].Out[k].Id == next.Neurons[j].In[m].Id {
						next.Neurons[j].In[m] = l.Neurons[i].Out[k]
					}
				}
			}
		}
	}
}

// ApplyBias creates and returns a bias synapse for each neuron in l
func (l *Layer) ApplyBias(weight synapse.WeightInitializer) []*synapse.Synapse {
	biases := make([]*synapse.Synapse, len(l.Neurons))
	for i := range l.Neurons {
		biases[i] = synapse.NewSynapse(weight())
		biases[i].IsBias = true
		l.Neurons[i].In = append(l.Neurons[i].In, biases[i])
	}
	return biases
}

func (l Layer) String() string {
	weights := make([][]float64, len(l.Neurons))
	for i, n := range l.Neurons {
		weights[i] = make([]float64, len(n.In))
		for j, s := range n.In {
			weights[i][j] = s.Weight
		}
	}
	return fmt.Sprintf("%+v", weights)
}
