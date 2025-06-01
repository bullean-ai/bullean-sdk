package activation

import (
	"math"
)

// Softmax is the softmax function
func Softmax(xx []float64) []float64 {
	out := make([]float64, len(xx))
	var sum float64
	maxx := Max(xx)
	for i, x := range xx {
		out[i] = math.Exp(x - maxx)
		sum += out[i]
	}
	for i := range out {
		out[i] /= sum
	}
	return out
}
