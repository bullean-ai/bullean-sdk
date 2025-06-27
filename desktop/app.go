package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bullean-ai/bullean-go/data"
	"github.com/bullean-ai/bullean-go/data/domain"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	wsClient domain.IClient
	Candles  map[string][]domain.Candle
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		Candles: make(map[string][]domain.Candle),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {

	quoteAsset := "XRP"
	baseAsset := "USDT"

	a.wsClient = data.NewClient(domain.ClientConfig{
		Version:   domain.V1,
		Name:      "Test",
		ApiKey:    "74918b1f-ba37-4c4e-a51a-9f26e9ceec2a",
		ApiSecret: "23a534d4c5cbcf86a9e77748e3f483",
		StreamReqMsg: domain.StreamReqMsg{
			TypeOf:      "subscription",
			History:     true,
			HistorySize: 3000,
			Subscriptions: []domain.Subscription{
				{
					Key:   "kline",
					Value: fmt.Sprintf("%s%s", quoteAsset, baseAsset),
				},
			},
		},
	})
	a.wsClient.OnReady(func(candles []domain.Candle) {
		for _, c := range candles {
			runtime.EventsEmit(ctx, "candles.init", candles)
			a.Candles[c.Symbol] = append(a.Candles[c.Symbol], c)
		}
	})
	a.wsClient.OnCandle(func(candles []domain.Candle) {

		for _, c := range candles {
			if c.Symbol == fmt.Sprintf("%s%s", quoteAsset, baseAsset) {
				candleBytes, err := json.Marshal(c)
				if err != nil {
					fmt.Println("InitCandles ERROR: ", err.Error())
				}

				runtime.EventsEmit(ctx, "candles.new", string(candleBytes))
			}
		}
	})
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) InitCandles(symbol string) string {
	_, isOk := a.Candles[symbol]
	if isOk {
		candleBytes, err := json.Marshal(a.Candles[symbol])
		if err != nil {
			fmt.Println("InitCandles ERROR: ", err.Error())
		}

		return string(candleBytes)
	}

	return ""
}
