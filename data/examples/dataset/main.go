package main

import (
	"context"
	"fmt"
	"github.com/bullean-ai/bullean-sdk/data"
	"github.com/bullean-ai/bullean-sdk/data/domain"
	"github.com/bullean-ai/bullean-sdk/data/neural/ffnn"
	"github.com/bullean-ai/bullean-sdk/data/neural/ffnn/solver"
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
			HistorySize: 10000,
			Subscriptions: []domain.Subscription{
				{
					Key:   "kline",
					Value: "BNBUSDT",
				},
			},
		},
	})
	var candless []domain.Candle
	var examples domain.Examples
	ranger := 100

	neural := ffnn.NewNeural(ffnn.DefaultConfig(ranger))

	client.OnReady(func(candles []domain.Candle) {

		dataset := data.NewDataSet(candles)

		dataset.CreatePolicy(domain.PolicyConfig{
			FeatName:    "feature_per_change",
			FeatType:    domain.FEAT_CLOSE_PERCENTAGE,
			PolicyRange: ranger,
		}, data.ClosePercentagePolicy)

		dataFrame := dataset.GetDataSet()

		for i := 0; i < len(dataFrame); i++ {
			label := []float64{}
			if dataFrame[i].Label == 1 {
				label = []float64{1, 0}
			} else if dataFrame[i].Label == 2 {
				label = []float64{0, 1}

			} else {
				label = []float64{0, 1}
			}
			examples = append(examples, domain.Example{
				Input:    dataFrame[i].Features,
				Response: label,
			})
		}
		candless = candles
		trainer := ffnn.NewTrainer(solver.NewSGD(0.005, 0.5, 1e-6, true), 1)
		//trainer := ffnn.NewTrainer(solver.NewAdam(0.001, 0, 0, 1e-15), 1)
		trainer.Train(neural, examples, examples, 100)

	})

	client.OnCandle(func(candles []domain.Candle) {
		for _, candle := range candles {
			if candle.Symbol == "BNBUSDT" {
				candless = append(candless, candle)
				dataset := data.NewDataSet(candless)

				dataset.CreatePolicy(domain.PolicyConfig{
					FeatName:    "feature_per_change",
					FeatType:    domain.FEAT_CLOSE_PERCENTAGE,
					PolicyRange: ranger,
				}, data.ClosePercentagePolicy)

				dataFrame := dataset.GetDataSet()

				for i := ranger; i < len(dataFrame); i++ {
					label := []float64{}
					if dataFrame[i].Label == 1 {
						label = []float64{1, 0}
					} else if dataFrame[i].Label == 2 {
						label = []float64{0, 1}

					} else {
						label = []float64{0, 1}
					}
					examples = append(examples, domain.Example{
						Input:    dataFrame[i].Features,
						Response: label,
					})
				}
				prediction := neural.Predict(examples[len(examples)-1].Input)
				fmt.Println(prediction)
			}
		}

	})

	data.GracefulExit(context.Background())
	return
}
