package ffnn

import (
	"fmt"
	"github.com/bullean-ai/bullean-go/neural_nets/domain"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/solver"
	"time"
)

// Trainer is a neurals network trainer
type Trainer interface {
	Train(n *FFNN, examples, validation domain.Examples, iterations int)
	FeedForward(n *FFNN, e domain.Example)
	BackPropagate(n *FFNN, e domain.Example, it int)
}

// OnlineTrainer is a basic, online network trainer
type OnlineTrainer struct {
	*internal
	solver     solver.Solver
	printer    *StatsPrinter
	verbosity  int
	prevProfit float64
	profit     float64
}

// NewTrainer creates a new trainer
func NewTrainer(solver solver.Solver, verbosity int) domain.ITrainer {
	return &OnlineTrainer{
		solver:    solver,
		printer:   NewStatsPrinter(),
		verbosity: verbosity,
	}
}

type internal struct {
	deltas [][]float64
}

func newTraining(layers []*layer.Layer) *internal {
	deltas := make([][]float64, len(layers))
	for i, l := range layers {
		deltas[i] = make([]float64, len(l.Neurons))
	}
	return &internal{
		deltas: deltas,
	}
}

// Train trains n
func (t *OnlineTrainer) Train(n interface{}, examples, validation domain.Examples, iterations int) (float64, domain.IModel) {
	var neural *FFNN

	switch n.(type) {
	case *FFNN:
		neural = n.(*FFNN)
	default:
		fmt.Println("Trainer can only train FFNN")
		return 0, nil
	}

	t.internal = newTraining(neural.Layers)

	train := make(domain.Examples, len(examples))
	copy(train, examples)

	t.printer.Init(neural)
	t.solver.Init(neural.NumWeights())
	accuracy := .0
	ts := time.Now()
	examples.Shuffle()
	for i := 1; i <= iterations; i++ {
		for j := 0; j < len(examples); j++ {
			t.FeedForward(neural, examples[j])
			t.BackPropagate(neural, examples[j], i)
		}
		if t.verbosity > 0 && i%t.verbosity == 0 && len(validation) > 0 {
			accuracy = t.printer.PrintProgress(neural, validation, time.Since(ts), i)
		}
	}
	return accuracy, neural
}

func (t *OnlineTrainer) learn(n *FFNN, e domain.Example, it int) {
	n.Forward(e.Input)
	t.calculateDeltas(n, e.Response)
	t.update(n, it)
}

func (t *OnlineTrainer) FeedForward(n *FFNN, e domain.Example) {
	n.Forward(e.Input)
}

func (t *OnlineTrainer) BackPropagate(n *FFNN, e domain.Example, it int) {
	t.calculateDeltas(n, e.Response)
	t.update(n, it)
}

func (t *OnlineTrainer) calculateDeltas(n *FFNN, ideal []float64) {
	for i, neuron := range n.Layers[len(n.Layers)-1].Neurons {
		t.deltas[len(n.Layers)-1][i] = GetLoss(n.Config.Loss).Df(
			neuron.Value,
			ideal[i],
			neuron.DActivate(neuron.Value))
	}

	for i := len(n.Layers) - 2; i >= 0; i-- {
		for j, neuron := range n.Layers[i].Neurons {
			var sum float64
			for k, s := range neuron.Out {
				sum += s.Weight * t.deltas[i+1][k]
			}
			t.deltas[i][j] = neuron.DActivate(neuron.Value) * sum
		}
	}
}

func (t *OnlineTrainer) update(n *FFNN, it int) {
	var idx int
	for i, l := range n.Layers {
		for j := range l.Neurons {
			for k := range l.Neurons[j].In {
				update := t.solver.Update(l.Neurons[j].In[k].Weight,
					t.deltas[i][j]*l.Neurons[j].In[k].In,
					it,
					idx)
				l.Neurons[j].In[k].Weight += update
				idx++
			}
		}
	}
}
