package ffnn

import (
	"fmt"
	"github.com/bullean-ai/bullean-go/neural_nets/domain"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/solver"
	"sync"
	"time"
)

// BatchTrainer implements parallelized batch training
type BatchTrainer struct {
	*internalb
	verbosity   int
	batchSize   int
	parallelism int
	solver      solver.Solver
	printer     *StatsPrinter
}

type internalb struct {
	deltas            [][][]float64
	partialDeltas     [][][][]float64
	accumulatedDeltas [][][]float64
	moments           [][][]float64
}

func newBatchTraining(layers []*layer.Layer, parallelism int) *internalb {
	deltas := make([][][]float64, parallelism)
	partialDeltas := make([][][][]float64, parallelism)
	accumulatedDeltas := make([][][]float64, len(layers))
	for w := 0; w < parallelism; w++ {
		deltas[w] = make([][]float64, len(layers))
		partialDeltas[w] = make([][][]float64, len(layers))

		for i, l := range layers {
			deltas[w][i] = make([]float64, len(l.Neurons))
			accumulatedDeltas[i] = make([][]float64, len(l.Neurons))
			partialDeltas[w][i] = make([][]float64, len(l.Neurons))
			for j, n := range l.Neurons {
				partialDeltas[w][i][j] = make([]float64, len(n.In))
				accumulatedDeltas[i][j] = make([]float64, len(n.In))
			}
		}
	}
	return &internalb{
		deltas:            deltas,
		partialDeltas:     partialDeltas,
		accumulatedDeltas: accumulatedDeltas,
	}
}

// NewBatchTrainer returns a BatchTrainer
func NewBatchTrainer(solver solver.Solver, verbosity, batchSize, parallelism int) *BatchTrainer {
	return &BatchTrainer{
		solver:      solver,
		verbosity:   verbosity,
		batchSize:   Iparam(batchSize, 1),
		parallelism: Iparam(parallelism, 1),
		printer:     NewStatsPrinter(),
	}
}

// Train trains n
func (t *BatchTrainer) Train(n interface{}, examples, validation domain.Examples, iterations int) (float64, domain.IModel) {
	var neural *FFNN

	switch n.(type) {
	case *FFNN:
		neural = n.(*FFNN)
	default:
		fmt.Println("Trainer can only train FFNN")
		return 0, nil
	}

	t.internalb = newBatchTraining(neural.Layers, t.parallelism)

	train := make(domain.Examples, len(examples))
	copy(train, examples)

	workCh := make(chan domain.Example, t.parallelism)
	nets := make([]*FFNN, t.parallelism)

	wg := sync.WaitGroup{}
	for i := 0; i < t.parallelism; i++ {
		nets[i] = NewFFNN(neural.Config)

		go func(id int, workCh <-chan domain.Example) {
			n := nets[id]
			for e := range workCh {
				n.Forward(e.Input)
				t.calculateDeltas(n, e.Response, id)
				wg.Done()
			}
		}(i, workCh)
	}

	t.printer.Init(neural)
	t.solver.Init(neural.NumWeights())

	ts := time.Now()
	for it := 1; it <= iterations; it++ {
		train.Shuffle()
		batches := train.SplitSize(t.batchSize)

		for _, b := range batches {
			currentWeights := neural.Weights()
			for _, n := range nets {
				n.ApplyWeights(currentWeights)
			}

			wg.Add(len(b))
			for _, item := range b {
				workCh <- item
			}
			wg.Wait()

			for _, wPD := range t.partialDeltas {
				for i, iPD := range wPD {
					iAD := t.accumulatedDeltas[i]
					for j, jPD := range iPD {
						jAD := iAD[j]
						for k, v := range jPD {
							jAD[k] += v
							jPD[k] = 0
						}
					}
				}
			}

			t.update(neural, it)
		}

		if t.verbosity > 0 && it%t.verbosity == 0 && len(validation) > 0 {
			t.printer.PrintProgress(neural, validation, time.Since(ts), it)
		}
	}
	return 0, neural
}

func (t *BatchTrainer) calculateDeltas(n *FFNN, ideal []float64, wid int) {
	loss := GetLoss(n.Config.Loss)
	deltas := t.deltas[wid]
	partialDeltas := t.partialDeltas[wid]
	lastDeltas := deltas[len(n.Layers)-1]

	for i, n := range n.Layers[len(n.Layers)-1].Neurons {
		lastDeltas[i] = loss.Df(
			n.Value,
			ideal[i],
			n.DActivate(n.Value))
	}

	for i := len(n.Layers) - 2; i >= 0; i-- {
		l := n.Layers[i]
		iD := deltas[i]
		nextD := deltas[i+1]
		for j, n := range l.Neurons {
			var sum float64
			for k, s := range n.Out {
				sum += s.Weight * nextD[k]
			}
			iD[j] = n.DActivate(n.Value) * sum
		}
	}

	for i, l := range n.Layers {
		iD := deltas[i]
		iPD := partialDeltas[i]
		for j, n := range l.Neurons {
			jD := iD[j]
			jPD := iPD[j]
			for k, s := range n.In {
				jPD[k] += jD * s.In
			}
		}
	}
}

func (t *BatchTrainer) update(n *FFNN, it int) {
	var idx int
	for i, l := range n.Layers {
		iAD := t.accumulatedDeltas[i]
		for j, n := range l.Neurons {
			jAD := iAD[j]
			for k, s := range n.In {
				update := t.solver.Update(s.Weight,
					jAD[k],
					it,
					idx)
				s.Weight += update
				jAD[k] = 0
				idx++
			}
		}
	}
}
