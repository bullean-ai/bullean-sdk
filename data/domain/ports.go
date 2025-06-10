package domain

type IClient interface {
	OnCandle(func([]Candle))
	OnReady(func([]Candle))
}

type IDataSet interface {
	CreatePolicy(PolicyConfig, func([]Candle) int)
	GetDataSet() []Data
	SerializeLabels()
}
