package domain

type IBinanceClient interface {
	Buy(BuyInfo)
	Sell(SellInfo)
}
