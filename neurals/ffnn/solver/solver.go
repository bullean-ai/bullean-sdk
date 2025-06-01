package solver

func Fparam(val, fallback float64) float64 {
	if val == 0.0 {
		return fallback
	}
	return val
}
