package indicators

import "github.com/bullean-ai/bullean-go/data/domain"

func MA(candles []domain.Candle, period int) (outs []float64) {
	for i := period; i < len(candles); i++ {
		sum := 0.
		for j := i - period; j < i; j++ {
			sum += candles[j].Close
		}
		outs = append(outs, sum/float64(period))
	}
	return
}
