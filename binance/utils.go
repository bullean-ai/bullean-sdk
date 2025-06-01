package binance

import "strconv"

func ToFloat(key interface{}) float64 {
	v, err := strconv.ParseFloat(key.(string), 64)
	if err != nil {
		return 0
	}
	return v
}
