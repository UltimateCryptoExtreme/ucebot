package main

import (
	"context"
	"strings"

	"github.com/gabu/go-cryptopia"
)

type Cryptopia struct {
	c *cryptopia.Client
}

func (c *Cryptopia) Name() string {
	return "Cryptopia"
}

func (c *Cryptopia) Init() {
	c.c = cryptopia.NewClient()
}

func (c *Cryptopia) CurrencyMarket(currency string) string {
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

func (c *Cryptopia) Price(market string) (Price, error) {
	market = strings.Replace(market, "-", "_", -1)
	p := strings.Split(market, "_")

	ctx := context.Background()
	data, err := c.c.GetMarket(ctx, market, 24)
	if err != nil {
		return Price{}, err
	}
	return Price{
		Source:     c.Name(),
		Quote:      p[0],
		Base:       p[1],
		Price:      data.LastPrice,
		Change24h:  &data.Change,
		Low:        data.Low,
		High:       data.High,
		BaseVolume: data.BaseVolume,
		LowestAsk:  data.AskPrice,
		HighestBid: data.BidPrice,
	}, nil
}
