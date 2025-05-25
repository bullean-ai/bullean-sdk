package main

import (
	"fmt"
	"github.com/bullean-ai/bullean-sdk/data"
	"github.com/bullean-ai/bullean-sdk/data/domain"
	"time"
)

func main() {

	now := time.Now()

	dataset := data.NewDataSet([]domain.Candle{
		{OpenTime: &now, Open: 100., High: 100., Low: 100., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 101., High: 101., Low: 101., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 102., High: 102., Low: 102., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 103., High: 103., Low: 103., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 104., High: 104., Low: 104., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 100., High: 100., Low: 100., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 98., High: 98., Low: 98., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 97., High: 97., Low: 97., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 95., High: 95., Low: 95., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 80., High: 80., Low: 80., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 82., High: 82., Low: 82., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 83., High: 83., Low: 83., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 84., High: 84., Low: 84., Close: 100., CloseTime: &now, Volume: 200.},
		{OpenTime: &now, Open: 85., High: 85., Low: 85., Close: 100., CloseTime: &now, Volume: 200.},
	})

	dataset.CreatePolicy(domain.PolicyConfig{
		FeatName:    "feature_open",
		FeatType:    domain.FEAT_OPEN,
		PolicyRange: 2,
	}, data.ClosePercentagePolicy)

	fmt.Println(dataset.GetDataSet())
	return
}
