package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bullean-ai/bullean-go/binance/domain"
)

type binanceFuturesClient struct {
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	client       *futures.Client
	lastOrder    *futures.CreateOrderResponse
	lastPosition int64
}

func NewBinanceFuturesClient(config domain.BinanceClientConfig) domain.IBinanceClient {
	var client = futures.NewClient(config.ApiKey, config.ApiSecret)
	account, _ := client.NewGetAccountService().Do(context.Background())
	if account == nil {
		fmt.Println("Account Info Error, check your api restrictions")
		return &binanceSpotClient{}

	}
	client.NewSetServerTimeService().Do(context.Background())
	return &binanceFuturesClient{
		ApiKey:       config.ApiKey,
		ApiSecret:    config.ApiSecret,
		client:       client,
		lastOrder:    nil,
		lastPosition: -1,
	}
}

func (b *binanceFuturesClient) Buy(req_dat domain.BuyInfo) {
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

		for _, bl := range account.Assets {
			if bl.Asset == req_dat.BaseAsset {
				balance = ToFloat(bl.AvailableBalance)
				fmt.Println(ToFloat(bl.AvailableBalance))
			}
		}
		var posType futures.SideType
		var posSide futures.PositionSideType
		switch req_dat.PositionSide {
		case 1:
			posType = futures.SideTypeBuy
			posSide = futures.PositionSideTypeLong
		case 2:
			posType = futures.SideTypeSell
			posSide = futures.PositionSideTypeShort
		default:
			return
		}
		quantity := (balance / req_dat.Price) * 99.9 / 100
		b.lastOrder, err = b.client.NewCreateOrderService().Symbol(fmt.Sprintf("%s%s", req_dat.QuoteAsset, req_dat.BaseAsset)).PositionSide(posSide).
			Side(posType).Type(futures.OrderTypeLimit).
			TimeInForce(futures.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%.2f", RoundDown(quantity, 2))).
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

func (b *binanceFuturesClient) Sell(req_dat domain.SellInfo) {
	var balance float64
	var err error
	account, _ := b.client.NewGetAccountService().Do(context.Background())
	if account != nil {
		for _, bl := range account.Assets {
			if bl.Asset == req_dat.QuoteAsset {
				balance = ToFloat(bl.AvailableBalance)
				f := ToFloat(bl.AvailableBalance)
				fmt.Println(fmt.Sprintf("%.5f", f))
			}
		}
		var posType futures.SideType
		var posSide futures.PositionSideType
		switch req_dat.PositionSide {
		case 1:
			posType = futures.SideTypeBuy
			posSide = futures.PositionSideTypeLong
		case 2:
			posType = futures.SideTypeSell
			posSide = futures.PositionSideTypeShort
		default:
			return
		}
		b.lastOrder, err = b.client.NewCreateOrderService().Symbol(fmt.Sprintf("%s%s", req_dat.QuoteAsset, req_dat.BaseAsset)).PositionSide(posSide).
			Side(posType).Type(futures.OrderTypeMarket).Quantity(fmt.Sprintf("%.2f", RoundDown(balance*req_dat.Price*99/100, 2))).Do(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		b.lastPosition = -1
	} else {
		fmt.Println("Account hatası")
	}
}
func (b *binanceFuturesClient) GetSymbolBalance(asset string) (balance float64) {
	account, _ := b.client.NewGetAccountService().Do(context.Background())
	if account != nil {
		for _, bl := range account.Assets {
			if bl.Asset == asset {
				balance = ToFloat(bl.AvailableBalance)
				break
			}
		}
	}

	return
}
