package datasources

import (
	"context"
	"errors"
	"net/http"

	"github.com/rmarku/ltp_api/internal/entities"
)

type DataSource interface {
	GetData(ctx context.Context, pair string) (*entities.LTP, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var ErrRequestFailed = errors.New("request failed")
