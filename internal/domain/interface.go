package domain

import "github.com/rmarku/ltp_api/internal/entities"

type LastTradePrice interface {
	GetLastTradePrices() (*entities.LTP, error)
}
