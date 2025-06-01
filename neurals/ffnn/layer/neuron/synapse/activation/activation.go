package activation

import "github.com/bullean-ai/bullean-sdk/neurals/domain"

// OutputActivation returns activation corresponding to prediction mode
func OutputActivation(c domain.Mode) domain.ActivationType {
	switch c {
	case domain.ModeMultiClass:
		return domain.ActivationSoftmax
	case domain.ModeRegression:
		return domain.ActivationLinear
	case domain.ModeBinary, domain.ModeMultiLabel:
		return domain.ActivationSigmoid
	}
	return domain.ActivationNone
}

// GetActivation returns the concrete activation given an ActivationType
func GetActivation(act domain.ActivationType) domain.Differentiable {
	switch act {
	case domain.ActivationSigmoid:
		return &Sigmoid{}
	case domain.ActivationTanh:
		return &Tanh{}
	case domain.ActivationReLU:
		return &ReLU{}
	case domain.ActivationLinear:
		return &Linear{}
	case domain.ActivationSoftmax:
		return &Linear{}
	}
	return &Linear{}
}

// Max is the largest element
func Max(xx []float64) float64 {
	max := xx[0]
	for _, x := range xx {
		if x > max {
			max = x
		}
	}
	return max
}

// ArgMax is the index of the largest element
func ArgMax(xx []float64) int {
	max, idx := xx[0], 0
	for i, x := range xx {
		if x > max {
			max, idx = xx[i], i
		}
	}
	return idx
}

// Normalize scales to (0,1)
func Normalize(xx []float64) {
	min, max := Min(xx), Max(xx)
	for i, x := range xx {
		xx[i] = (x - min) / (max - min)
	}
}

// Min is the smallest element
func Min(xx []float64) float64 {
	min := xx[0]
	for _, x := range xx {
		if x < min {
			min = x
		}
	}
	return min
}
