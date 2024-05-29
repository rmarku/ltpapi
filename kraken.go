package main

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type LTP struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type krakenResponse struct {
	Result map[string]Result `json:"result"`
	Error  []interface{}     `json:"error"`
}

type Result struct {
	Data []string `json:"c"`
}

const tickRoot = "https://api.kraken.com/0/public/Ticker"

func getTicket(pair string) (*LTP, error) {
	var response krakenResponse
	req, _ := http.NewRequest("GET", tickRoot, nil)
	q := req.URL.Query()
	q.Add("pair", pair)
	req.URL.RawQuery = q.Encode()
	slog.Debug("Built URL", "url", req.URL.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Error) > 0 {
		return nil, errors.New("Kraken returned error")
	}

	amount, err := strconv.ParseFloat(response.Result[pair].Data[0], 64)
	if err != nil {
		return nil, err
	}

	return &LTP{Pair: pair, Amount: amount}, nil
}
