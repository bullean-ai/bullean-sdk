package main

import (
	"context"
	"fmt"
	"github.com/bullean-ai/bullean-sdk/data"
	"github.com/bullean-ai/bullean-sdk/data/domain"
	ffnnDomain "github.com/bullean-ai/bullean-sdk/neurals/domain"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/solver"
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
	var examples ffnnDomain.Examples
	ranger := 300

	neural := ffnn.NewNeural(ffnn.DefaultConfig(ranger))

	client.OnReady(func(candles []domain.Candle) {

		dataset := data.NewDataSet(candles)

		dataset.CreatePolicy(domain.PolicyConfig{
			FeatName:    "feature_per_change",
			FeatType:    domain.FEAT_CLOSE_PERCENTAGE,
			PolicyRange: 100,
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
			examples = append(examples, ffnnDomain.Example{
				Input:    dataFrame[i].Features,
				Response: label,
			})
		}
		candless = candles
		trainer := ffnn.NewTrainer(solver.NewAdam(0.001, 0, 0, 1e-12), 1)
		//trainer := ffnn.NewBatchTrainer(solver.NewSGD(0.0005, 0.1, 0, true), 1, ranger, 12)
		trainer.Train(neural, examples, examples, 1000)

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
					examples = append(examples, ffnnDomain.Example{
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
