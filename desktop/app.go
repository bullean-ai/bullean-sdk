package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bullean-ai/bullean-go/data"
	"github.com/bullean-ai/bullean-go/data/domain"
	"github.com/bullean-ai/bullean-go/neural_nets"
	ffnnDomain "github.com/bullean-ai/bullean-go/neural_nets/domain"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/layer/neuron/synapse"
	"github.com/bullean-ai/bullean-go/neural_nets/ffnn/solver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"math"
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
			HistorySize: 1200,
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
			a.Candles[c.Symbol] = append(a.Candles[c.Symbol], c)
		}
		candleBytes, err := json.Marshal(candles)
		if err != nil {
			fmt.Println("InitCandles ERROR: ", err.Error())
		}
		runtime.EventsEmit(ctx, "candles.init", string(candleBytes))
		runtime.EventsEmit(ctx, "candles.done", true)
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

func (a *App) GetPredictions(symbol string) {
	fmt.Println(symbol)
	var examples ffnnDomain.Examples
	inputLen := 700
	ranger := 20
	iterations := 50
	var lr float64 = 0.004
	var model1 *ffnn.FFNN
	//var err error

	model1 = ffnn.NewFFNN(&ffnnDomain.Config{
		Inputs:     inputLen + 1,
		Layout:     []int{30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 3},
		Activation: ffnnDomain.ActivationSoftmax,
		Mode:       ffnnDomain.ModeMultiClass,
		Weight:     synapse.NewNormal(1e-20, 1e-20),
		Bias:       true,
	})
	evaluator := neural_nets.NewEvaluator([]ffnnDomain.Neural{
		{
			Model:      model1,
			Trainer:    ffnn.NewBatchTrainer(solver.NewAdam(lr, 0, 0, 1e-12), 1, 100, 12),
			Iterations: iterations,
		},
	})

	dataset := data.NewDataSet(a.Candles[symbol], inputLen)

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
			Time:     *dataFrame[i].Time,
		})
	}

	evaluator.Train(examples, examples)

	var predictions []struct {
		Time       string `json:"time"`
		Prediction int8   `json:"prediction"`
	}

	for i := 0; i < len(examples); i++ {
		var prediction int8

		pred := evaluator.Predict(examples[i].Input)
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

		predictions = append(predictions, struct {
			Time       string `json:"time"`
			Prediction int8   `json:"prediction"`
		}{
			Time:       examples[i].Time.Format("2006-01-02T15:04:05Z"),
			Prediction: prediction + 5,
		})
	}
	pBytes, err := json.Marshal(predictions)
	if err != nil {
		fmt.Println("InitCandles ERROR: ", err.Error())
	}
	runtime.EventsEmit(a.ctx, "candles.prediction", string(pBytes))
}
