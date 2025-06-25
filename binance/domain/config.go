package domain

type BinanceClientConfig struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

type BuyInfo struct {
	BaseAsset    string   `json:"BaseAsset"`
	QuoteAsset   string   `json:"QuoteAsset"`
	Amount       float64  `json:"amount"`
	Price        float64  `json:"price"`
	StopLoss     *float64 `json:"stop_loss"`
	TakeProfit   *float64 `json:"take_profit"`
	PositionSide int      `json:"position_side"` //TODO: For Futures implementation
}

type SellInfo struct {
	BaseAsset    string  `json:"BaseAsset"`
	QuoteAsset   string  `json:"QuoteAsset"`
	Price        float64 `json:"price"`
	PositionSide int     `json:"position_side"`
}
