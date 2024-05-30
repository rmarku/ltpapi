package domain

import (
	"errors"

	"github.com/rmarku/ltp_api/internal/entities"
)

type LastTradePrice interface {
	GetLastTradePrices(pair string) (*entities.LTP, error)
	UpdatePrices() error
	GetPairs() []string
}

var ErrPriceNotFound = errors.New("price not found")
