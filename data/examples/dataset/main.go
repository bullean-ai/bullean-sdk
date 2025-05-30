package main

import (
	"context"
	"fmt"
	"github.com/bullean-ai/bullean-sdk/data"
	"github.com/bullean-ai/bullean-sdk/data/domain"
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
			HistorySize: 100,
			Subscriptions: []domain.Subscription{
				{
					Key:   "kline",
					Value: "BNBUSDT",
				},
			},
		},
	})
	var candless []domain.Candle

	client.OnReady(func(candles []domain.Candle) {
		dataset := data.NewDataSet(candles)

		dataset.CreatePolicy(domain.PolicyConfig{
			FeatName:    "feature_open",
			FeatType:    domain.FEAT_OPEN,
			PolicyRange: 50,
		}, data.ClosePercentagePolicy)

		dataFrame := dataset.GetDataSet()
		for _, dat := range dataFrame {
			fmt.Println(dat.Features[0], " : ", dat.Features[len(dat.Features)-1], " : ", dat.Label)
		}
		candless = candles
	})

	client.OnCandle(func(candles []domain.Candle) {
		for _, candle := range candles {
			candless = append(candless, candle)
		}
		fmt.Println(candles)
	})

	data.GracefulExit(context.Background())
	return
}
