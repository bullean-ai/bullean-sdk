package synapse

// Synapse is an edge between neurons
type Synapse struct {
	Id      int
	Weight  float64
	In, Out float64
	IsBias  bool
}

// NewSynapse returns a synapse with the specified initialized weight
func NewSynapse(weight float64) *Synapse {
	return &Synapse{Weight: weight}
}

func (s *Synapse) Fire(value float64) {
	s.In = value
	s.Out = s.In * s.Weight
}
