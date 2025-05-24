package domain

type IClient interface {
	OnCandle(func(Candle))
	OnReady() []Candle
}
