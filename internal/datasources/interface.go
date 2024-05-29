package datasources

import (
	"context"
	"errors"
	"net/http"
)

type DataSource interface {
	GetData(ctx context.Context, pair string) (*LTP, error)
}

type LTP struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var ErrRequestFailed = errors.New("request failed")
