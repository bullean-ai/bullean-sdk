package indicators

import "github.com/bullean-ai/bullean-go/data/domain"

func EMA(candles []domain.Candle, period int) (outs []float64) {
	for i := 1; i < len(candles); i++ {
		k := 2.0 / float64(1+period)
		val := (candles[i].Close * k) + (candles[i-1].Close * (1 - k))
		outs = append(outs, val)
	}
	return
}
