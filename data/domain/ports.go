package domain

type IClient interface {
	OnCandle(func(Candle))
	OnReady() []Candle
}

type IDataSet interface {
	CreatePolicy(PolicyConfig, func([]Candle) int)
	GetDataSet() []Data
}
