package livecoin

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const BASE_URL = "https://api.livecoin.net/"

type Livecoin struct {
	BaseUrl *url.URL
}

func New() *Livecoin {
	l := &Livecoin{}
	l.BaseUrl, _ = url.Parse(BASE_URL)
	return l
}

func (l *Livecoin) newReq(ctx context.Context, method, relUrl string) (*http.Request, error) {
	rel, err := url.Parse(relUrl)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	u := l.BaseUrl.ResolveReference(rel)
	req, err = http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (l *Livecoin) do(req *http.Request, v interface{}) error {
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
