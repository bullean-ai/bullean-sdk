package binance

import (
	"math"
	"strconv"
)

func ToFloat(key interface{}) float64 {
	v, err := strconv.ParseFloat(key.(string), 64)
	if err != nil {
		return 0
	}
	return v
}

func RoundDown(val float64, precision int) float64 {
	return math.Floor(val*(math.Pow10(precision))) / math.Pow10(precision)
}
