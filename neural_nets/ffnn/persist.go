package ffnn

import (
	"encoding/json"
	"fmt"
	"github.com/bullean-ai/bullean-go/neural_nets/domain"
	"os"
)

// Dump is a neurals network dump
type Dump struct {
	Config  *domain.Config
	Weights [][][]float64
}

// ApplyWeights sets the weights from a three-dimensional slice
func (n *FFNN) ApplyWeights(weights [][][]float64) {
	for i, l := range n.Layers {
		for j := range l.Neurons {
			for k := range l.Neurons[j].In {
				n.Layers[i].Neurons[j].In[k].Weight = weights[i][j][k]
			}
		}
	}
}

// Weights returns all weights in sequence
func (n FFNN) Weights() [][][]float64 {
	weights := make([][][]float64, len(n.Layers))
	for i, l := range n.Layers {
		weights[i] = make([][]float64, len(l.Neurons))
		for j, n := range l.Neurons {
			weights[i][j] = make([]float64, len(n.In))
			for k, in := range n.In {
				weights[i][j][k] = in.Weight
			}
		}
	}
	return weights
}

// Dump generates a network dump
func (n *FFNN) Dump() *Dump {
	return &Dump{
		Config:  n.Config,
		Weights: n.Weights(),
	}
}

// FromDump restores a FFNN from a dump
func FromDump(dump *Dump) *FFNN {
	n := NewFFNN(dump.Config)
	n.ApplyWeights(dump.Weights)

	return n
}

// Marshal marshals to JSON from network
func (n *FFNN) Marshal() ([]byte, error) {
	return json.Marshal(n.Dump())
}

func (n *FFNN) SaveModel(path string) error {
	bytes, err := n.Marshal()
	if err != nil {
		fmt.Println("Error marshaling network:", err.Error())
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err.Error())
	}

	return err
}

func LoadModel(path string) (ffnn *FFNN, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(path)
	if err != nil {
		fmt.Println("Error writing to file:", err.Error())
	}
	ffnn, err = Unmarshal(bytes)
	if err != nil {
		fmt.Println("Error marshaling network:", err.Error())
		return nil, err
	}
	return ffnn, err
}

// Unmarshal restores network from a JSON blob
func Unmarshal(bytes []byte) (*FFNN, error) {
	var dump Dump
	if err := json.Unmarshal(bytes, &dump); err != nil {
		return nil, err
	}
	return FromDump(&dump), nil
}
