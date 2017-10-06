package livecoin

import (
	"context"
)

type Ticker struct {
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Average float64 `json:"vwap"`
	Volume  float64 `json:"volume"`
	Last    float64 `json:"last"`
	Buy     float64 `json:"best_bid"`
	Sell    float64 `json:"best_ask"`
}

func (l *Livecoin) TickerAll(ctx context.Context) ([]Ticker, error) {
	url := "exchange/ticker"
	req, err := l.newReq(ctx, "GET", url)
	if err != nil {
		return nil, err
	}
	var tickers []Ticker
	err = l.do(req, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}

func (l *Livecoin) Ticker(ctx context.Context, market string) (Ticker, error) {
	url := "exchange/ticker?currencyPair=" + market
	req, err := l.newReq(ctx, "GET", url)
	if err != nil {
		return Ticker{}, err
	}
	var ticker Ticker
	err = l.do(req, &ticker)
	if err != nil {
		return Ticker{}, err
	}
	return ticker, nil
}
