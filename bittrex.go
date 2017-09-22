package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toorop/go-bittrex"
)

type Bittrex struct {
	c *bittrex.Bittrex
}

func (b *Bittrex) Name() string {
	return "Bittrex"
}

func (b *Bittrex) Init() {
	b.c = bittrex.New("", "")
}

func (b *Bittrex) CurrencyMarket(currency string) string {
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

func (b *Bittrex) Price(market string) (Price, error) {
	pair := strings.Split(market, "-")
	market = fmt.Sprintf("%s-%s", pair[1], pair[0])

	ss, err := b.c.GetMarketSummary(market)
	if err != nil {
		return Price{}, err
	}
	if len(ss) == 0 {
		return Price{}, errors.New("Market not found.")
	}
	data := ss[0]
	change := ((data.Last - data.PrevDay) / data.PrevDay) * 100
	return Price{
		Source:     b.Name(),
		Quote:      pair[0],
		Base:       pair[1],
		Price:      data.Last,
		Change24h:  &change,
		Low:        data.Low,
		High:       data.High,
		BaseVolume: data.BaseVolume,
		LowestAsk:  data.Ask,
		HighestBid: data.Bid,
	}, nil
}
