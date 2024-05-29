package domain

import "github.com/rmarku/ltp_api/internal/entities"

type LastTradePrice interface {
	GetLastTradePrices(pair string) (*entities.LTP, error)
	UpdatePrices() error
}
