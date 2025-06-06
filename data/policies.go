package data

import "github.com/bullean-ai/bullean-go/data/domain"

// ClosePercentagePolicy is a default close price percentage change policy
func ClosePercentagePolicy(candles []domain.Candle) int {
	perChange := 0.
	for i, _ := range candles {
		if i == 0 {
			continue
		}
		perChange += ((candles[i].Close - candles[i-1].Close) / candles[i-1].Close) * 100
	}
	if perChange > .4 {
		return 1
	} else if perChange < .4 && perChange >= .4 {
		return 0
	} else {
		return 2
	}
}
