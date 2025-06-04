package ffnn

import (
	"encoding/json"
	"fmt"
	"github.com/bullean-ai/bullean-go/neural_nets/domain"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer"
	synapse "github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer/neuron/synapse"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer/neuron/synapse/activation"
	"os"
)

// FFNN is a neurals network
type FFNN struct {
	Layers []*layer.Layer
	Biases [][]*synapse.Synapse
	Config *domain.Config
}

// NewFFNN returns a new neurals network
func NewFFNN(c *domain.Config) *FFNN {

	if c.Weight == nil {
		c.Weight = synapse.NewUniform(0.5, 0)
	}
	if c.Activation == domain.ActivationNone {
		c.Activation = domain.ActivationSigmoid
	}
	if c.Loss == domain.LossNone {
		switch c.Mode {
		case domain.ModeMultiClass, domain.ModeMultiLabel:
			c.Loss = domain.LossCrossEntropy
		case domain.ModeBinary:
			c.Loss = domain.LossBinaryCrossEntropy
		default:
			c.Loss = domain.LossMeanSquared
		}
	}

	layers := initializeLayers(c)

	var biases [][]*synapse.Synapse
	if c.Bias {
		biases = make([][]*synapse.Synapse, len(layers))
		for i := 0; i < len(layers); i++ {
			if c.Mode == domain.ModeRegression && i == len(layers)-1 {
				continue
			}
			biases[i] = layers[i].ApplyBias(c.Weight)
		}
	}

	return &FFNN{
		Layers: layers,
		Biases: biases,
		Config: c,
	}
}

func initializeLayers(c *domain.Config) []*layer.Layer {
	layers := make([]*layer.Layer, len(c.Layout))
	for i := range layers {
		act := c.Activation
		if i == (len(layers)-1) && c.Mode != domain.ModeDefault {
			act = activation.OutputActivation(c.Mode)
		}
		layers[i] = layer.NewLayer(c.Layout[i], act)
	}

	for i := 0; i < len(layers)-1; i++ {
		layers[i].Connect(layers[i+1], c.Weight)
	}

	for _, neuron := range layers[0].Neurons {
		neuron.In = make([]*synapse.Synapse, c.Inputs)
		for i := range neuron.In {
			neuron.In[i] = synapse.NewSynapse(c.Weight())
		}
	}

	return layers
}

func ConnectPreparedFFNN(neural *FFNN) {

	for i := 0; i < len(neural.Layers)-1; i++ {
		neural.Layers[i].ConnectPrepared(neural.Layers[i+1])
	}

	return
}

func (n *FFNN) Fire() {
	for _, b := range n.Biases {
		for _, s := range b {
			s.Fire(1)
		}
	}
	for _, l := range n.Layers {
		l.Fire()
	}
}

// Forward computes a forward pass
func (n *FFNN) Forward(input []float64) error {
	if len(input) != n.Config.Inputs {
		return fmt.Errorf("Invalid input dimension - expected: %d got: %d", n.Config.Inputs, len(input))
	}
	for _, n := range n.Layers[0].Neurons {
		for i := 0; i < len(input); i++ {
			n.In[i].Fire(input[i])
		}
	}
	n.Fire()
	return nil
}

// Predict computes a forward pass and returns a prediction
func (n *FFNN) Predict(input []float64) []float64 {
	n.Forward(input)

	outLayer := n.Layers[len(n.Layers)-1]
	out := make([]float64, len(outLayer.Neurons))
	for i, neuron := range outLayer.Neurons {
		out[i] = neuron.Value
	}
	return out
}

// NumWeights returns the number of weights in the network
func (n *FFNN) NumWeights() (num int) {
	for _, l := range n.Layers {
		for _, n := range l.Neurons {
			num += len(n.In)
		}
	}
	return
}

func (n *FFNN) String() string {
	var s string
	for _, l := range n.Layers {
		s = fmt.Sprintf("%s\n%s", s, l)
	}
	return s
}

func (n *FFNN) Save(directory string) (err error) {
	var bytes []byte
	var file *os.File

	file, err = os.Open(directory)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	bytes, err = json.Marshal(n)

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return
}
