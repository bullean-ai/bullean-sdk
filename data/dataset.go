package data

import "github.com/bullean-ai/bullean-sdk/data/domain"

type dataSet struct {
	Candles     []domain.Candle `json:"candles"`
	Features    []domain.Data   `json:"features"`
	PolicyRange int             `json:"policy_range"`
}

func NewDataSet(policy_range int, candles []domain.Candle) domain.IDataSet {
	return &dataSet{
		Candles:     candles,
		PolicyRange: policy_range,
	}
}

// CreatePolicy Policy is created by the callback function.
// 0)Hold 1)Buy 2)Sell
func (d *dataSet) CreatePolicy(config domain.PolicyConfig, policy func([]domain.Candle) int) {
	for i := len(d.Candles) - d.PolicyRange - 1; i > 0; i-- {
		signal := policy(d.Candles[i : i+d.PolicyRange])
		data := domain.Data{
			Name:     config.FeatName,
			Features: d.GetFeatureValues(d.Candles[i:i+d.PolicyRange], config.FeatType),
			Label:    signal,
		}
		d.Features = append(d.Features, data)
	}
}

func (d *dataSet) GetDataSet() (data []domain.Data) {

	data = d.Features

	return
}

func (d *dataSet) GetFeatureValues(candles []domain.Candle, feat_type domain.FeatureType) (data []float64) {
	for i, candle := range candles {
		switch feat_type {
		case domain.FEAT_OPEN:
			data = append(data, candle.Open)
		case domain.FEAT_HIGH:
			data = append(data, candle.High)
		case domain.FEAT_LOW:
			data = append(data, candle.Low)
		case domain.FEAT_CLOSE:
			data = append(data, candle.Close)
		case domain.FEAT_CLOSE_PERCENTAGE:
			if i == 0 {
				data = append(data, 0)
			} else {
				data = append(data, ((candles[i].Close-candles[i-1].Close)/candles[i-1].Close)*100)

			}
		}
	}
	return
}

// ClosePercentagePolicy is a default policy
func ClosePercentagePolicy(candles []domain.Candle) int {
	perChange := 0.
	for i, _ := range candles {
		if i == 0 {
			continue
		}
		perChange += ((candles[i].Open - candles[i-1].Open) / candles[i-1].Open) * 100
	}
	if perChange > 0.01 {
		return 1
	} else if perChange < 0.01 && perChange > -.01 {
		return 0
	} else {
		return -1
	}
}
