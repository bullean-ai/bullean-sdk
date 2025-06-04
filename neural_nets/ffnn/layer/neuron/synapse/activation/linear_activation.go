package activation

import "github.com/bullean-ai/bullean-go/neural_nets/domain"

// Linear is a linear activator
type Linear struct{}

func NewLinearActivation() domain.Differentiable {
	return &Linear{}
}

// F is the identity function
func (a *Linear) F(x float64) float64 { return x }

// Df is constant
func (a *Linear) Df(x float64) float64 { return 1 }
