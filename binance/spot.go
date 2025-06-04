package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bullean-ai/bullean-go/binance/domain"
	"math"
)

type binanceClient struct {
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	client       *binance.Client
	lastOrder    *binance.CreateOrderResponse
	lastPosition int64
}

func NewBinanceClient(config domain.BinanceClientConfig) domain.IBinanceClient {
	var client = binance.NewClient(config.ApiKey, config.ApiSecret)
	account, _ := client.NewGetAccountService().Do(context.Background())
	if account == nil {
		fmt.Println("Account Info Error, check your api restrictions")
		return nil
	}
	client.NewSetServerTimeService().Do(context.Background())
	return &binanceClient{
		ApiKey:       config.ApiKey,
		ApiSecret:    config.ApiSecret,
		client:       client,
		lastOrder:    nil,
		lastPosition: -1,
	}
}

func (b *binanceClient) Buy(req_dat domain.BuyInfo) {
	var balance float64
	var err error

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
			Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
			TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%.2f", RoundDown((balance/req_dat.Price)*99.9/100, 2))).
			Price(fmt.Sprintf("%.2f", req_dat.Price)).Do(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("Account hatası")
	}

	b.lastPosition = 1

}

func (b *binanceClient) Sell(req_dat domain.SellInfo) {
	var balance float64
	var err error
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
			Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
			TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%.2f", RoundDown(balance*99.9/100, 2))).
			Price(fmt.Sprintf("%.2f", req_dat.Price)).Do(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		b.lastPosition = -1
	} else {
		fmt.Println("Account hatası")
	}

}
func RoundDown(val float64, precision int) float64 {
	return math.Floor(val*(math.Pow10(precision))) / math.Pow10(precision)
}
