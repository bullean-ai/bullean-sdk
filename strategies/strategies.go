package strategies

import (
	"fmt"
	binanceDomain "github.com/bullean-ai/bullean-go/binance/domain"
	"github.com/bullean-ai/bullean-go/data/domain"
	domain2 "github.com/bullean-ai/bullean-go/strategies/domain"
)

type Strategy struct {
	BaseAsset           string
	QuoteAssets         []string
	BinanceClient       binanceDomain.IBinanceClient
	Candles             map[string][]domain.Candle
	CandleLimit         int
	LastLongEnterPrice  float64
	LastLongClosePrice  float64
	LastShortEnterPrice float64
	LastShortClosePrice float64
}

func NewStrategy(base_asset string, quote_assets []string, candle_limit int, client binanceDomain.IBinanceClient) *Strategy {
	return &Strategy{
		BaseAsset:     base_asset,
		QuoteAssets:   quote_assets,
		BinanceClient: client,
		CandleLimit:   candle_limit,
		Candles:       map[string][]domain.Candle{},
	}
}

func (s *Strategy) Next(candle []domain.Candle) {
	for _, c := range candle {
		if len(s.Candles[c.Symbol]) >= s.CandleLimit {
			s.Candles[c.Symbol] = s.Candles[c.Symbol][1:] // Remove the oldest candle if limit exceeded
		}
		if _, exists := s.Candles[c.Symbol]; !exists {
			s.Candles[c.Symbol] = []domain.Candle{}
		}
		s.Candles[c.Symbol] = append(s.Candles[c.Symbol], c) // Add the new candle
	}
}

func (s *Strategy) Evaluate(long_condition func(last_long_enter_price, last_long_close_price float64) domain2.PositionType, short_condition func(last_short_price, last_short_close_price float64) domain2.PositionType) {
	for _, quote := range s.QuoteAssets {
		if _, exists := s.Candles[fmt.Sprintf("%s%s", quote, s.BaseAsset)]; exists {
			if long_condition(s.LastLongEnterPrice, s.LastLongClosePrice) == domain2.POS_BUY {
				s.LongEnter(quote)
			} else if long_condition(s.LastLongEnterPrice, s.LastLongClosePrice) == domain2.POS_SELL {
				s.LongClose(quote)
			}
			if short_condition(s.LastShortEnterPrice, s.LastShortClosePrice) == domain2.POS_BUY {
				s.ShortEnter(quote)
			} else if short_condition(s.LastShortEnterPrice, s.LastShortClosePrice) == domain2.POS_SELL {
				s.ShortClose(quote)
			}
		}
	}
}

func (s *Strategy) LongEnter(quote_asset string) {
	candleSeries := s.Candles[fmt.Sprintf("%s%s", quote_asset, s.BaseAsset)]
	s.LastLongEnterPrice = candleSeries[len(candleSeries)-1].Close
	balance := s.BinanceClient.GetSymbolBalance(quote_asset)
	s.BinanceClient.Buy(binanceDomain.BuyInfo{
		BaseAsset:    s.BaseAsset,
		QuoteAsset:   quote_asset,
		Amount:       balance,
		Price:        s.LastLongEnterPrice,
		PositionSide: 1,
	})
}

func (s *Strategy) ShortEnter(quote_asset string) {
	candleSeries := s.Candles[fmt.Sprintf("%s%s", quote_asset, s.BaseAsset)]
	s.LastShortEnterPrice = candleSeries[len(candleSeries)-1].Close
	balance := s.BinanceClient.GetSymbolBalance(quote_asset)
	s.BinanceClient.Buy(binanceDomain.BuyInfo{
		BaseAsset:    s.BaseAsset,
		QuoteAsset:   quote_asset,
		Amount:       balance,
		Price:        s.LastShortEnterPrice,
		PositionSide: 2,
	})

}

func (s *Strategy) LongClose(quote_asset string) {
	candleSeries := s.Candles[fmt.Sprintf("%s%s", quote_asset, s.BaseAsset)]
	s.LastLongClosePrice = candleSeries[len(candleSeries)-1].Close
	s.BinanceClient.Sell(binanceDomain.SellInfo{
		BaseAsset:    s.BaseAsset,
		QuoteAsset:   quote_asset, // Assuming selling to the first quote asset
		Price:        s.LastLongClosePrice,
		PositionSide: 1,
	})
}

func (s *Strategy) ShortClose(quote_asset string) {
	candleSeries := s.Candles[fmt.Sprintf("%s%s", quote_asset, s.BaseAsset)]
	s.LastShortClosePrice = candleSeries[len(candleSeries)-1].Close
	s.BinanceClient.Sell(binanceDomain.SellInfo{
		BaseAsset:    s.BaseAsset,
		QuoteAsset:   quote_asset, // Assuming selling to the first quote asset
		Price:        s.LastShortClosePrice,
		PositionSide: 2,
	})
}
