package domain

func PercentageChange(val1, val2 float64) float64 {
	if val1 == 0 {
		val1 = 1
	}
	return ((val2 - val1) / val1) * 100
}
