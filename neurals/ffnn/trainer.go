package ffnn

import (
	"github.com/bullean-ai/bullean-sdk/neurals/domain"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/layer"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/solver"
	"time"
)

// Trainer is a neurals network trainer
type Trainer interface {
	Train(n *Neural, examples, validation domain.Examples, iterations int)
	FeedForward(n *Neural, e domain.Example)
	BackPropagate(n *Neural, e domain.Example, it int)
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
func NewTrainer(solver solver.Solver, verbosity int) *OnlineTrainer {
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
func (t *OnlineTrainer) Train(n *Neural, examples, validation domain.Examples, iterations int) float64 {
	t.internal = newTraining(n.Layers)

	train := make(domain.Examples, len(examples))
	copy(train, examples)

	t.printer.Init(n)
	t.solver.Init(n.NumWeights())
	accuracy := .0
	ts := time.Now()
	examples.Shuffle()
	for i := 1; i <= iterations; i++ {
		for j := 0; j < len(examples); j++ {
			t.FeedForward(n, examples[j])
			t.BackPropagate(n, examples[j], i)
		}
		if t.verbosity > 0 && i%t.verbosity == 0 && len(validation) > 0 {
			accuracy = t.printer.PrintProgress(n, validation, time.Since(ts), i)
		}
	}
	return accuracy
}

func (t *OnlineTrainer) learn(n *Neural, e domain.Example, it int) {
	n.Forward(e.Input)
	t.calculateDeltas(n, e.Response)
	t.update(n, it)
}

func (t *OnlineTrainer) FeedForward(n *Neural, e domain.Example) {
	n.Forward(e.Input)
}

func (t *OnlineTrainer) BackPropagate(n *Neural, e domain.Example, it int) {
	t.calculateDeltas(n, e.Response)
	t.update(n, it)
}

func (t *OnlineTrainer) calculateDeltas(n *Neural, ideal []float64) {
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

func (t *OnlineTrainer) update(n *Neural, it int) {
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
