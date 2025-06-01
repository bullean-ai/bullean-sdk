package solver

// Solver implements an update rule for training a NN
type Solver interface {
	Init(size int)
	Update(value, gradient float64, iteration, idx int) float64
}

// SGD is stochastic gradient descent with nesterov/momentum
type SGD struct {
	lr       float64
	decay    float64
	momentum float64
	nesterov bool
	moments  []float64
}

// NewSGD returns a new SGD solver
func NewSGD(lr, momentum, decay float64, nesterov bool) *SGD {
	return &SGD{
		lr:       Fparam(lr, 0.01),
		decay:    decay,
		momentum: momentum,
		nesterov: nesterov,
	}
}

// Init initializes vectors using number of weights in network
func (o *SGD) Init(size int) {
	o.moments = make([]float64, size)
}

// Update returns the update for a given weight
func (o *SGD) Update(value, gradient float64, iteration, idx int) float64 {
	lr := o.lr / (1 + o.decay*float64(iteration))

	o.moments[idx] = o.momentum*o.moments[idx] - lr*gradient

	if o.nesterov {
		o.moments[idx] = o.momentum*o.moments[idx] - lr*gradient
	}

	return o.moments[idx]
}
