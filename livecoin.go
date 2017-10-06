package main

import (
	"context"
	"strings"

	"github.com/ultimatecryptoextreme/go-livecoin"
)

type Livecoin struct {
	c *livecoin.Livecoin
}

func (l *Livecoin) Name() string {
	return "Livecoin"
}

func (l *Livecoin) Init() {
	l.c = livecoin.New()
}

func (l *Livecoin) CurrencyMarket(currency string) string {
	if strings.Contains(currency, "-") {
		return currency
	}
	switch currency {
	case "BTC":
		return currency + "-USD"
	default:
		return currency + "-BTC"
	}
}

func (l *Livecoin) Price(market string) (Price, error) {
	market = strings.Replace(market, "-", "/", -1)
	p := strings.Split(market, "/")
	market = strings.ToUpper(market)

	ctx := context.Background()
	ticker, err := l.c.Ticker(ctx, market)
	if err != nil {
		return Price{}, err
	}
	return Price{
		Source:     l.Name(),
		Quote:      p[0],
		Base:       p[1],
		Price:      ticker.Last,
		Low:        ticker.Low,
		High:       ticker.High,
		BaseVolume: ticker.Volume * ticker.Average,
		LowestAsk:  ticker.Sell,
		HighestBid: ticker.Buy,
	}, nil
}
