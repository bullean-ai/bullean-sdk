package domain

type IBinanceClient interface {
	Buy(BuyInfo)
	Sell(SellInfo)
	GetSymbolBalance(string) float64
}
