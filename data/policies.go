package data

import "github.com/bullean-ai/bullean-sdk/data/domain"

// ClosePercentagePolicy is a default close price percentage change policy
func ClosePercentagePolicy(candles []domain.Candle) int {
	perChange := 0.
	for i, _ := range candles {
		if i == 0 {
			continue
		}
		perChange += ((candles[i].Open - candles[i-1].Open) / candles[i-1].Open) * 100
	}
	if perChange > 0.01 {
		return 1
	} else if perChange < 0.01 && perChange > -.01 {
		return 0
	} else {
		return -1
	}
}
