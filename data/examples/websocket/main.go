package main

import (
	"context"
	"github.com/bullean-ai/bullean-go/data"
	"github.com/bullean-ai/bullean-go/data/domain"
)

func main() {

	client := data.NewClient(domain.ClientConfig{
		Version:   domain.V1,
		Name:      "Test",
		ApiKey:    "74918b1f-ba37-4c4e-a51a-9f26e9ceec2a",
		ApiSecret: "23a534d4c5cbcf86a9e77748e3f483",
		StreamReqMsg: domain.StreamReqMsg{
			TypeOf:      "subscription",
			History:     true,
			HistorySize: 30000,
			Subscriptions: []domain.Subscription{
				{
					Key:   "kline",
					Value: "BTCUSDT",
				},
				{
					Key:   "kline",
					Value: "BNBUSDT",
				},
			},
		},
	})

	client.OnCandle(func(candles []domain.Candle) {
		println("Candle received:", candles)
	})

	data.GracefulExit(context.Background())

}
