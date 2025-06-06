package neural_nets

import (
	"fmt"
	"github.com/bullean-ai/bullean-go/neural_nets/domain"
	"math"
	"sync"
)

type Evaluator struct {
	Neurals []domain.Neural
}

func NewEvaluator(neurals []domain.Neural) *Evaluator {
	return &Evaluator{
		Neurals: neurals,
	}
}

func (n *Evaluator) Train(train_data domain.Examples, validate_date domain.Examples) {
	wg := &sync.WaitGroup{}

	tAmount := len(train_data) / len(n.Neurals)
	vAmount := len(validate_date) / len(n.Neurals)
	for i := 0; i < len(n.Neurals); i++ {
		t_truncated := train_data[i*tAmount : (i+1)*tAmount]
		v_truncated := validate_date[i*vAmount : (i+1)*vAmount]
		wg.Add(1)
		go func(neural domain.Neural, t_data domain.Examples, v_data domain.Examples, index int) {
			_, n.Neurals[index].Model = neural.Trainer.Train(neural.Model, t_data, v_data, n.Neurals[index].Iterations)
			wg.Done()
		}(n.Neurals[i], t_truncated, v_truncated, i)
	}

	wg.Wait()

	var deciderSamples domain.Examples

	for _, sample := range train_data {
		var s domain.Example
		predictions := make([]float64, len(n.Neurals))
		for i := 0; i < len(n.Neurals); i++ {
			pred := n.Neurals[i].Model.Predict(sample.Input)
			predictions = append(predictions, math.Round(pred[0]))
		}
		s.Input = predictions
		s.Response = sample.Response
		deciderSamples = append(deciderSamples, s)
	}

	// Evaluate the last neural with all data
	n.Neurals[len(n.Neurals)-1].Trainer.Train(n.Neurals[len(n.Neurals)-1].Model, deciderSamples, deciderSamples, n.Neurals[len(n.Neurals)-1].Iterations)

	return
}

func (n *Evaluator) Predict(input []float64) []float64 {
	predictions := make([][]float64, len(n.Neurals))
	var lastChoice []float64
	buys := 0
	sells := 0
	for i := 0; i < len(n.Neurals); i++ {
		predictions = append(predictions, n.Neurals[i].Model.Predict(input))
		pred := n.Neurals[i].Model.Predict(input)
		buy := math.Round(pred[0])
		if buy == 1 {
			buys += 1
		} else {
			sells += 1
		}

	}

	fmt.Println(predictions)
	if buys > sells {
		lastChoice = []float64{1, 0}
	} else {
		lastChoice = []float64{0, 1}
	}

	return lastChoice
}
