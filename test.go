package main

import (
	"context"
	"fmt"
	"github.com/bullean-ai/bullean-go/binance"
	binanceDomain "github.com/bullean-ai/bullean-go/binance/domain"
	"github.com/bullean-ai/bullean-go/data"
	"github.com/bullean-ai/bullean-go/data/domain"
	"github.com/bullean-ai/bullean-go/neural_nets"
	ffnnDomain "github.com/bullean-ai/bullean-go/neural_nets/domain"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/solver"
	"math"
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
			HistorySize: 40000,
			Subscriptions: []domain.Subscription{
				{
					Key:   "kline",
					Value: "BNBUSDT",
				},
			},
		},
	})
	binanceClient := binance.NewBinanceSpotClient(binanceDomain.BinanceClientConfig{
		ApiKey:    "xsCifNydCadBpp4gmL3TJGmNlNWt6nkKTfdhOmEmFqM8WqPqyXBAoDk97YviHorI",
		ApiSecret: "6KgbU2UNZ8rdjhkOLNf1QWJ3c941sTbkaouf0TJ0OSAiFkIqpq4n8MWUKeQCE7st",
	})
	var candless []domain.Candle
	var examples ffnnDomain.Examples
	isReady := false
	ranger := 100

	evaluator := neural_nets.NewEvaluator([]ffnnDomain.Neural{
		{
			Model:   ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(ranger)),
			Trainer: ffnn.NewTrainer(solver.NewAdam(0.002, 0, 0, 1e-12), 1),
		},
		{
			Model:   ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(ranger)),
			Trainer: ffnn.NewTrainer(solver.NewAdam(0.002, 0, 0, 1e-12), 1),
		},
		{
			Model:   ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(ranger)),
			Trainer: ffnn.NewTrainer(solver.NewAdam(0.002, 0, 0, 1e-12), 1),
		},
	})

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
			examples = append(examples, ffnnDomain.Example{
				Input:    dataFrame[i].Features,
				Response: label,
			})
		}
		candless = candles
		//trainer := ffnn.NewBatchTrainer(solver.NewSGD(0.0005, 0.1, 0, true), 1, ranger, 12)

		evaluator.Train(examples, examples, 20)
		isReady = true

	})

	lastprediction := 0
	client.OnCandle(func(candles []domain.Candle) {
		var prediction int
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
				pred := evaluator.Predict(examples[len(examples)-1].Input)
				buy := math.Round(pred[0])
				if buy == 1 {
					prediction = 1
				} else {
					prediction = -1
				}

				fmt.Println(pred)
				if isReady == false {
					continue
				}

				if prediction == 1 && lastprediction == 1 {
					binanceClient.Buy(binanceDomain.BuyInfo{
						Price:      candle.Close,
						QuoteAsset: "BNB",
						BaseAsset:  "USDC",
					})
				} else if prediction == -1 && lastprediction == -1 {
					binanceClient.Sell(binanceDomain.SellInfo{
						Price:      candle.Close,
						QuoteAsset: "BNB",
						BaseAsset:  "USDC",
					})
				}

				lastprediction = prediction
			}
		}

	})

	data.GracefulExit(context.Background())
	return
}
