package main

import (
	"context"
	"fmt"
	"github.com/bullean-ai/bullean-go/binance"
	binanceDomain "github.com/bullean-ai/bullean-go/binance/domain"
	"github.com/bullean-ai/bullean-go/data"
	"github.com/bullean-ai/bullean-go/data/domain"
	"github.com/bullean-ai/bullean-go/indicators"
	"github.com/bullean-ai/bullean-go/neural_nets"
	ffnnDomain "github.com/bullean-ai/bullean-go/neural_nets/domain"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/solver"
	"github.com/bullean-ai/bullean-go/strategies"
	domain2 "github.com/bullean-ai/bullean-go/strategies/domain"
	"math"
)

func main() {

	quoteAsset := "XRP"
	baseAsset := "USDT"

	client := data.NewClient(domain.ClientConfig{
		Version:   domain.V1,
		Name:      "Test",
		ApiKey:    "74918b1f-ba37-4c4e-a51a-9f26e9ceec2a",
		ApiSecret: "23a534d4c5cbcf86a9e77748e3f483",
		StreamReqMsg: domain.StreamReqMsg{
			TypeOf:      "subscription",
			History:     true,
			HistorySize: 4100,
			Subscriptions: []domain.Subscription{
				{
					Key:   "kline",
					Value: fmt.Sprintf("%s%s", quoteAsset, baseAsset),
				},
			},
		},
	})
	binanceClient := binance.NewBinanceFuturesClient(binanceDomain.BinanceClientConfig{
		ApiKey:    "xsCifNydCadBpp4gmL3TJGmNlNWt6nkKTfdhOmEmFqM8WqPqyXBAoDk97YviHorI",
		ApiSecret: "6KgbU2UNZ8rdjhkOLNf1QWJ3c941sTbkaouf0TJ0OSAiFkIqpq4n8MWUKeQCE7st",
	})
	strategy := strategies.NewStrategy(baseAsset, []string{quoteAsset}, 40, binanceClient)

	var candless []domain.Candle
	var examples ffnnDomain.Examples
	var willTrain = true
	isReady := false
	inputLen := 2000
	ranger := 10
	iterations := 50
	lr := 0.004
	var model1 *ffnn.FFNN
	//var err error

	model1 = ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(inputLen))
	/*
		model1, err = ffnn.LoadModel("./model1.json")
		if err != nil {
			fmt.Println(err.Error())
		}
			model1, err = ffnn.LoadModel("./model1.json")
				if err != nil {
				 		model1 = ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(inputLen))
					willTrain = true
				}
				model2, err = ffnn.LoadModel("./model2.json")
				if err != nil {
					model2 = ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(inputLen))
					willTrain = true
				}
	*/
	evaluator := neural_nets.NewEvaluator([]ffnnDomain.Neural{
		{
			Model:      model1,
			Trainer:    ffnn.NewBatchTrainer(solver.NewAdam(lr, 0, 0, 1e-12), 1, 100, 12),
			Iterations: iterations,
		},
	})

	client.OnReady(func(candles []domain.Candle) {

		dataset := data.NewDataSet(candles, inputLen)

		dataset.CreatePolicy(domain.PolicyConfig{
			FeatName:    "feature_per_change",
			FeatType:    domain.FEAT_CLOSE_PERCENTAGE,
			PolicyRange: ranger,
		}, data.MAPercentagePolicy)
		dataset.SerializeLabels()

		dataFrame := dataset.GetDataSet()
		for _, dat := range dataFrame {
			fmt.Println(dat.Label)
		}

		for i := 0; i < len(dataFrame); i++ {
			label := []float64{}
			if dataFrame[i].Label == 1 {
				label = []float64{1, 0, 0}
			} else if dataFrame[i].Label == 2 {
				label = []float64{0, 1, 0}

			} else {
				label = []float64{0, 0, 1}
			}
			examples = append(examples, ffnnDomain.Example{
				Input:    dataFrame[i].Features,
				Response: label,
			})
		}
		candless = candles
		//trainer := ffnn.NewBatchTrainer(solver.NewSGD(0.0005, 0.1, 0, true), 1, ranger, 12)
		if willTrain == true {
			evaluator.Train(examples, examples)
			model1.SaveModel("./model1.json")

		}
		isReady = true

	})

	lastprediction := 0
	//isTrainingEnd := true
	client.OnCandle(func(candles []domain.Candle) {
		var prediction int
		for _, candle := range candles {
			if candle.Symbol == "XRPUSDT" {
				candless = candless[1:]
				candless = append(candless, candle)
				if isReady == false {
					continue
				}
				dataset := data.NewDataSet(candless, inputLen)

				dataset.CreatePolicy(domain.PolicyConfig{
					FeatName:    "feature_per_change",
					FeatType:    domain.FEAT_CLOSE_PERCENTAGE,
					PolicyRange: ranger,
				}, data.MAPercentagePolicy)
				dataset.SerializeLabels()
				dataFrame := dataset.GetDataSet()

				for i := 0; i < len(dataFrame); i++ {
					label := []float64{}
					if dataFrame[i].Label == 1 {
						label = []float64{1, 0, 0}
					} else if dataFrame[i].Label == 2 {
						label = []float64{0, 1, 0}

					} else {
						label = []float64{0, 0, 1}
					}
					examples = append(examples, ffnnDomain.Example{
						Input:    dataFrame[i].Features,
						Response: label,
					})
				}
				/*
					go func() {
						if isTrainingEnd {
							isTrainingEnd = false
							model2 := ffnn.NewFFNN(ffnnDomain.DefaultFFNNConfig(ranger))
							newEvaluator := neural_nets.NewEvaluator([]ffnnDomain.Neural{
								{
									Model:      model2,
									Trainer:    ffnn.NewBatchTrainer(solver.NewAdam(lr, 0, 0, 1e-12), 1, 100, 12),
									Iterations: iterations,
								},
							})
							newEvaluator.Train(examples, examples)
							//model1.SaveModel("./model1.json")
							model2.SaveModel("./model2.json")
							*model1 = *model2
							newEvaluator = neural_nets.NewEvaluator([]ffnnDomain.Neural{
								{
									Model:      model1,
									Trainer:    ffnn.NewBatchTrainer(solver.NewAdam(lr, 0, 0, 1e-12), 1, 100, 12),
									Iterations: iterations,
								},
							})
							evaluator = newEvaluator
							isTrainingEnd = true
						}
					}()
				*/
				pred := evaluator.Predict(examples[len(examples)-1].Input)
				fmt.Println(pred, len(candless[len(candless)-16:]))
				buy := math.Round(pred[0])
				sell := math.Round(pred[1])
				hold := math.Round(pred[2])
				if buy == 1 {
					prediction = 1
				} else if sell == 1 {
					prediction = -1
				} else if hold == 1 {
					prediction = 0
				}

				emaOuts := indicators.EMA(candless[len(candless)-16:], 6)
				changedir := 0
				if domain2.PercentageChange(emaOuts[1], emaOuts[len(emaOuts)-1]) > 0 {
					changedir = 1
				} else if domain2.PercentageChange(emaOuts[1], emaOuts[len(emaOuts)-1]) < 0 {
					changedir = -1
				}
				strategy.Next(candles)
				strategy.Evaluate(func(lastLongEnterPrice, lastLongClosePrice float64) domain2.PositionType { // Long Enter
					fmt.Println("LONG: ", lastLongEnterPrice, lastLongClosePrice, domain2.PercentageChange(lastLongEnterPrice, candle.Close))
					if prediction == 1 && lastprediction == 1 {
						return domain2.POS_BUY
					} else if (prediction == 1 && lastprediction == 1) || changedir == -1 {
						return domain2.POS_SELL
					} else {
						return domain2.POS_HOLD
					}

				}, func(lastShortEnterPrice, lastShortClosePrice float64) domain2.PositionType { // Short Enter
					if lastShortEnterPrice == 0 {
						lastShortEnterPrice = candle.Close
					}
					fmt.Println("SHORT: ", lastShortEnterPrice, lastShortClosePrice, domain2.PercentageChange(lastShortEnterPrice, candle.Close))
					if prediction == -1 && lastprediction == -1 {
						return domain2.POS_BUY
					} else if (prediction == -1 && lastprediction == -1) || changedir == 1 {
						return domain2.POS_SELL
					} else {
						return domain2.POS_HOLD
					}
				})

				lastprediction = prediction
			}
		}

	})

	data.GracefulExit(context.Background())
	return
}
