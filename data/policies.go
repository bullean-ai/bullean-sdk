package data

import (
	"github.com/bullean-ai/bullean-go/data/domain"
	"github.com/bullean-ai/bullean-go/indicators"
	domain2 "github.com/bullean-ai/bullean-go/strategies/domain"
)

// ClosePercentagePolicy is a default close price percentage change policy
func ClosePercentagePolicy(candles []domain.Candle) int {
	perChange := 0.
	for i, _ := range candles {
		if i == 0 {
			continue
		}
		perChange += ((candles[i].Close - candles[i-1].Close) / candles[i-1].Close) * 100
	}
	if perChange > .3 {
		return 1
	} else if perChange < .3 && perChange >= 0 {
		return 0
	} else {
		return 2
	}
}

// MAPercentagePolicy is a default close price percentage change policy
func MAPercentagePolicy(candles []domain.Candle) int {
	ema := indicators.MA(candles, 5)
	if domain2.PercentageChange(ema[0], ema[len(ema)-1]) >= 0.3 {
		return 1
	} else if domain2.PercentageChange(ema[0], ema[len(ema)-1]) < 0.3 && domain2.PercentageChange(ema[0], ema[len(ema)-1]) >= 0 {
		return 0
	} else {
		return 2
	}
}
