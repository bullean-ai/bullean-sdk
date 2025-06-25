package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/bullean-ai/bullean-go/binance/domain"
)

type binanceSpotClient struct {
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	client       *binance.Client
	lastOrder    *binance.CreateOrderResponse
	lastPosition int64
}

func NewBinanceSpotClient(config domain.BinanceClientConfig) domain.IBinanceClient {
	var client = binance.NewClient(config.ApiKey, config.ApiSecret)
	account, _ := client.NewGetAccountService().Do(context.Background())
	if account == nil {
		fmt.Println("Account Info Error, check your api restrictions")
		return &binanceSpotClient{}
	}
	client.NewSetServerTimeService().Do(context.Background())
	return &binanceSpotClient{
		ApiKey:       config.ApiKey,
		ApiSecret:    config.ApiSecret,
		client:       client,
		lastOrder:    nil,
		lastPosition: -1,
	}

}

func (b *binanceSpotClient) Buy(req_dat domain.BuyInfo) {
	var balance float64
	var err error
	if b.client == nil {
		return
	}
	// Cancel Order
	if b.lastOrder != nil {
		_, err = b.client.NewCancelOrderService().Symbol(fmt.Sprintf("%s%s", req_dat.QuoteAsset, req_dat.BaseAsset)).
			OrderID(b.lastOrder.OrderID).Do(context.Background())
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	account, _ := b.client.NewGetAccountService().Do(context.Background())
	if account != nil {
		for _, bl := range account.Balances {
			if bl.Asset == req_dat.BaseAsset {
				balance = ToFloat(bl.Free)
				fmt.Println(ToFloat(bl.Free), ToFloat(bl.Locked))
			}
		}
		b.lastOrder, err = b.client.NewCreateOrderService().Symbol(fmt.Sprintf("%s%s", req_dat.QuoteAsset, req_dat.BaseAsset)).
			Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).Quantity(fmt.Sprintf("%.3f", RoundDown((balance/req_dat.Price)*99.9/100, 3))).
			Do(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("Account hatası")
	}

	b.lastPosition = 1
}

func (b *binanceSpotClient) Sell(req_dat domain.SellInfo) {
	var balance float64
	var err error
	if b.client == nil {
		return
	}
	account, _ := b.client.NewGetAccountService().Do(context.Background())
	if account != nil {
		for _, bl := range account.Balances {
			if bl.Asset == req_dat.QuoteAsset {
				balance = ToFloat(bl.Free)
				f := ToFloat(bl.Free)
				fmt.Println(fmt.Sprintf("%.5f", f), ToFloat(bl.Locked))
			}
		}
		b.lastOrder, err = b.client.NewCreateOrderService().Symbol(fmt.Sprintf("%s%s", req_dat.QuoteAsset, req_dat.BaseAsset)).
			Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).Quantity(fmt.Sprintf("%.3f", RoundDown(balance*99.9/100, 3))).
			Do(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		b.lastPosition = -1
	} else {
		fmt.Println("Account hatası")
	}
}

func (b *binanceSpotClient) GetSymbolBalance(asset string) (balance float64) {
	account, _ := b.client.NewGetAccountService().Do(context.Background())
	if account != nil {
		for _, bl := range account.Balances {
			if bl.Asset == asset {
				balance = ToFloat(bl.Free)
				break
			}
		}
	}

	return
}
