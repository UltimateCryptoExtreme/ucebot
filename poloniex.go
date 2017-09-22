package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jyap808/go-poloniex"
)

type Poloniex struct {
	c *poloniex.Poloniex
}

func (p *Poloniex) Name() string {
	return "Poloniex"
}

func (p *Poloniex) Init() {
	p.c = poloniex.New("", "")
}

func (p *Poloniex) CurrencyMarket(currency string) string {
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

func (p *Poloniex) Price(market string) (Price, error) {
	pair := strings.Split(market, "-")
	market = fmt.Sprintf("%s_%s", pair[1], pair[0])

	t, err := p.c.GetTickers()
	if err != nil {
		return Price{}, err
	}
	data, ok := t[market]
	if !ok {
		return Price{}, errors.New("Market not found.")
	}
	change := data.PercentChange * 100
	return Price{
		Source:     p.Name(),
		Quote:      pair[0],
		Base:       pair[1],
		Price:      data.Last,
		Change24h:  &change,
		Low:        data.Low24Hr,
		High:       data.High24Hr,
		BaseVolume: data.BaseVolume,
		LowestAsk:  data.LowestAsk,
		HighestBid: data.HighestBid,
	}, nil
}
