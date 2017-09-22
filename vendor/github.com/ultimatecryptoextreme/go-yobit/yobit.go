package yobit

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const BASE_URL = "https://yobit.net/api/3/"

type Yobit struct {
	BaseUrl *url.URL
}

func New() *Yobit {
	y := &Yobit{}
	y.BaseUrl, _ = url.Parse(BASE_URL)
	return y
}

func (y *Yobit) newReq(ctx context.Context, method, relUrl string) (*http.Request, error) {
	rel, err := url.Parse(relUrl)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	u := y.BaseUrl.ResolveReference(rel)
	req, err = http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (y *Yobit) do(req *http.Request, v interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}
