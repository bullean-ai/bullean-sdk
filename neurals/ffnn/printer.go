package ffnn

import (
	"fmt"
	"github.com/bullean-ai/bullean-sdk/neurals/domain"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/layer/neuron/synapse/activation"
	"os"
	"text/tabwriter"
	"time"
)

// StatsPrinter prints training progress
type StatsPrinter struct {
	w *tabwriter.Writer
}

// NewStatsPrinter creates a StatsPrinter
func NewStatsPrinter() *StatsPrinter {
	return &StatsPrinter{tabwriter.NewWriter(os.Stdout, 16, 0, 3, ' ', 0)}
}

// Init initializes printer
func (p *StatsPrinter) Init(n *Neural) {
	fmt.Fprintf(p.w, "Epochs\tElapsed\tLoss (%s)\t", n.Config.Loss)
	if n.Config.Mode == domain.ModeMultiClass {
		fmt.Fprintf(p.w, "Accuracy\t\n---\t---\t---\t---\t\n")
	} else {
		fmt.Fprintf(p.w, "\n---\t---\t---\t\n")
	}
}

// PrintProgress prints the current state of training
func (p *StatsPrinter) PrintProgress(n *Neural, validation domain.Examples, elapsed time.Duration, iteration int) float64 {
	fmt.Fprintf(p.w, "%d\t%s\t%.4f\t%s\n",
		iteration,
		elapsed.String(),
		CrossValidate(n, validation),
		FormatAccuracy(n, validation))
	p.w.Flush()
	return Accuracy(n, validation)
}

func FormatAccuracy(n *Neural, validation domain.Examples) string {
	if n.Config.Mode == domain.ModeMultiClass {
		return fmt.Sprintf("%.2f\t", Accuracy(n, validation))
	}
	return ""
}

func Accuracy(n *Neural, validation domain.Examples) float64 {
	correct := 0
	for _, e := range validation {
		est := n.Predict(e.Input)
		if activation.ArgMax(e.Response) == activation.ArgMax(est) {
			correct++
		}
	}
	return float64(correct) / float64(len(validation))
}

func CrossValidate(n *Neural, validation domain.Examples) float64 {
	predictions, responses := make([][]float64, len(validation)), make([][]float64, len(validation))
	for i := 0; i < len(validation); i++ {
		predictions[i] = n.Predict(validation[i].Input)
		responses[i] = validation[i].Response
	}

	return GetLoss(n.Config.Loss).F(predictions, responses)
}
