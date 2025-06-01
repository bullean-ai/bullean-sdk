package ffnn

import (
	"math"
	"strings"
)

// Mean of xx
func Mean(xx []float64) float64 {
	var sum float64
	for _, x := range xx {
		sum += x
	}
	return sum / float64(len(xx))
}

// Variance of xx
func Variance(xx []float64) float64 {
	if len(xx) == 1 {
		return 0.0
	}
	m := Mean(xx)

	var variance float64
	for _, x := range xx {
		variance += math.Pow((x - m), 2)
	}

	return variance / float64(len(xx)-1)
}

// StandardDeviation of xx
func StandardDeviation(xx []float64) float64 {
	return math.Sqrt(Variance(xx))
}

// Standardize (z-score) shifts distribution to μ=0 σ=1
func Standardize(xx []float64) {
	m := Mean(xx)
	s := StandardDeviation(xx)

	if s == 0 {
		s = 1
	}

	for i, x := range xx {
		xx[i] = (x - m) / s
	}
}

// Sgn is signum
func Sgn(x float64) float64 {
	switch {
	case x < 0:
		return -1.0
	case x > 0:
		return 1.0
	}
	return 0
}

// Sum is sum
func Sum(xx []float64) (sum float64) {
	for _, x := range xx {
		sum += x
	}
	return
}

// Round to nearest integer
func Round(x float64) float64 {
	return math.Floor(x + .5)
}

// Dot product
func Dot(xx, yy []float64) float64 {
	var p float64
	for i := range xx {
		p += xx[i] * yy[i]
	}
	return p
}

func Iparam(val, fallback int) int {
	if val == 0 {
		return fallback
	}
	return val
}

// CheckStringIfContains check a string if contains given param
func CheckStringIfContains(input_text string, search_text string) bool {
	CheckContains := strings.Contains(input_text, search_text)
	return CheckContains
}

// Sma Simple Moving Average.
func Sma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1
		sum += value

		if i >= period {
			sum -= values[i-period]
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

func MaxValue(values []float64) (result float64, index int) {
	result = math.MinInt64
	for i := 0; i < len(values); i++ {
		if values[i] > result {
			result = values[i]
			index = i

		}
	}
	return
}

func MinValue(values []float64) (result float64, index int) {
	result = math.MaxInt64
	for i := 0; i < len(values); i++ {
		if values[i] < result {
			result = values[i]
			index = i
		}
	}
	return
}
