package activation

import (
	"github.com/bullean-ai/bullean-sdk/data/domain"
	"math"
)

// Tanh is a hyperbolic activator
type Tanh struct{}

func NewTanhActivation() domain.Differentiable {
	return &Tanh{}
}

// F is Tanh(x)
func (a *Tanh) F(x float64) float64 { return (1 - math.Exp(-2*x)) / (1 + math.Exp(-2*x)) }

// Df is Tanh'(y), where y = Tanh(x)
func (a *Tanh) Df(y float64) float64 { return 1 - math.Pow(y, 2) }
