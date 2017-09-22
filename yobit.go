package main

import (
	"context"
	"strings"

	"github.com/ultimatecryptoextreme/go-yobit"
)

type Yobit struct {
	c *yobit.Yobit
}

func (y *Yobit) Name() string {
	return "Yobit"
}

func (y *Yobit) Init() {
	y.c = yobit.New()
}

func (y *Yobit) CurrencyMarket(currency string) string {
	if strings.Contains(currency, "-") {
		return currency
	}
	switch currency {
	case "BTC":
		return currency + "-USDT"
	default:
		return currency + "-BTC"
	}
}

func (y *Yobit) Price(market string) (Price, error) {
	market = strings.Replace(market, "-", "_", -1)
	p := strings.Split(market, "_")
	market = strings.ToLower(market)

	ctx := context.Background()
	tickers, err := y.c.Ticker(ctx, []string{market})
	if err != nil {
		return Price{}, err
	}
	data := tickers[market]
	return Price{
		Source:     y.Name(),
		Quote:      p[0],
		Base:       p[1],
		Price:      data.Last,
		Low:        data.Low,
		High:       data.High,
		BaseVolume: data.BaseVolume,
		LowestAsk:  data.Sell,
		HighestBid: data.Buy,
	}, nil
}
