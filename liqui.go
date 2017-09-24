package main

import (
	"context"
	"net/url"
	"strings"

	"github.com/ultimatecryptoextreme/go-yobit"
)

type Liqui struct {
	c *yobit.Yobit
}

func (l *Liqui) Name() string {
	return "Liqui"
}

func (l *Liqui) Init() {
	l.c = yobit.New()
	l.c.BaseUrl, _ = url.Parse("https://api.liqui.io/api/3/")
}

func (l *Liqui) CurrencyMarket(currency string) string {
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

func (l *Liqui) Price(market string) (Price, error) {
	market = strings.Replace(market, "-", "_", -1)
	p := strings.Split(market, "_")
	market = strings.ToLower(market)

	ctx := context.Background()
	tickers, err := l.c.Ticker(ctx, []string{market})
	if err != nil {
		return Price{}, err
	}
	data := tickers[market]
	return Price{
		Source:     l.Name(),
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
