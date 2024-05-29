package datasources

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type DataSourceKraken struct {
	client HTTPClient
	uri    string
}

var ErrKrakenRequest = errors.New("kraken request failed")

// Type check.
var _ DataSource = new(DataSourceKraken)

type krakenResponse struct {
	Result map[string]result `json:"result"`
	Error  []any             `json:"error"`
}

type result struct {
	Data []string `json:"c"`
}

func NewKraken(uri string) *DataSourceKraken {
	return &DataSourceKraken{uri: uri, client: &http.Client{}}
}

func (k *DataSourceKraken) GetData(ctx context.Context, pair string) (*LTP, error) {
	var response krakenResponse

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, k.uri, nil)
	q := req.URL.Query()

	q.Add("pair", pair)
	req.URL.RawQuery = q.Encode()
	slog.Debug("Built URL", "url", req.URL.String())

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, errors.Join(err, ErrRequestFailed)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(err, ErrRequestFailed)
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, errors.Join(err, ErrRequestFailed)
	}

	if len(response.Error) > 0 {
		return nil, errors.Join(ErrKrakenRequest, ErrRequestFailed)
	}

	amount, err := strconv.ParseFloat(response.Result[pair].Data[0], 64)
	if err != nil {
		return nil, errors.Join(err, ErrRequestFailed)
	}

	return &LTP{Pair: pair, Amount: amount}, nil
}
