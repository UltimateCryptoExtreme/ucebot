package yobit

import (
	"context"
	"strings"
)

type Ticker struct {
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Average     float64 `json:"avg"`
	BaseVolume  float64 `json:"vol"`
	QuoteVolume float64 `json:"vol_cur"`
	Last        float64 `json:"last"`
	Buy         float64 `json:"buy"`
	Sell        float64 `json:"sell"`
}

func (y *Yobit) Ticker(ctx context.Context, markets []string) (map[string]Ticker, error) {
	url := "ticker/" + strings.Join(markets, "-")
	req, err := y.newReq(ctx, "GET", url)
	if err != nil {
		return nil, err
	}
	var tickers map[string]Ticker
	err = y.do(req, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}
