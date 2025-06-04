package solver

import (
	"math"
)

// Adam is an Adam solver
type Adam struct {
	lr      float64
	beta    float64
	beta2   float64
	epsilon float64

	v, m []float64
}

// NewAdam returns a new Adam solver
func NewAdam(lr, beta, beta2, epsilon float64) *Adam {
	return &Adam{
		lr:      Fparam(lr, 0.001),
		beta:    Fparam(beta, 0.9),
		beta2:   Fparam(beta2, 0.999),
		epsilon: Fparam(epsilon, 1e-8),
	}
}

// Init initializes vectors using number of weights in network
func (o *Adam) Init(size int) {
	o.v, o.m = make([]float64, size), make([]float64, size)
}

// Update returns the update for a given weight
func (o *Adam) Update(value, gradient float64, t, idx int) float64 {
	lrt := o.lr * (math.Sqrt(1.0 - math.Pow(o.beta2, float64(t)))) /
		(1.0 - math.Pow(o.beta, float64(t)))
	o.m[idx] = o.beta*o.m[idx] + (1.0-o.beta)*gradient
	o.v[idx] = o.beta2*o.v[idx] + (1.0-o.beta2)*math.Pow(gradient, 2.0)

	return -lrt * (o.m[idx] / (math.Sqrt(o.v[idx]) + o.epsilon))
}
