package data

import "github.com/bullean-ai/bullean-go/data/domain"

type dataSet struct {
	Candles  []domain.Candle `json:"candles"`
	Features []domain.Data   `json:"features"`
	InputLen int             `json:"input_len"`
}

func NewDataSet(candles []domain.Candle, input_len int) domain.IDataSet {
	return &dataSet{
		Candles:  candles,
		InputLen: input_len,
	}
}

// CreatePolicy Policy is created by the callback function.
// 0)Hold 1)Buy 2)Sell
func (d *dataSet) CreatePolicy(config domain.PolicyConfig, policy func([]domain.Candle) int) {
	for i := len(d.Candles) - config.PolicyRange; i >= d.InputLen+1; i-- {
		var data domain.Data
		signal := policy(d.Candles[i : i+config.PolicyRange])
		if d.InputLen > 0 {
			data = domain.Data{
				Name:     config.FeatName,
				Features: d.GetFeatureValues(d.Candles[i-d.InputLen:i], config.FeatType),
				Label:    signal,
			}
		} else {
			data = domain.Data{
				Name:     config.FeatName,
				Features: d.GetFeatureValues(d.Candles[i:i+config.PolicyRange], config.FeatType),
				Label:    signal,
			}
		}

		d.Features = append([]domain.Data{data}, d.Features...)
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
