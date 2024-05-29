package domain

import (
	"context"

	"github.com/rmarku/ltp_api/internal/datasources"
	"github.com/rmarku/ltp_api/internal/entities"
)

type LastTradePriceImpl struct {
	source datasources.DataSource
	pairs  []string
}

var _ LastTradePrice = new(LastTradePriceImpl)

func NewLastTradePrice(source datasources.DataSource) *LastTradePriceImpl {

	return &LastTradePriceImpl{
		pairs:  []string{"BTC/USD", "BTC/CHF", "BTC/EUR"},
		source: source,
	}
}

func (l *LastTradePriceImpl) GetLastTradePrices() (*entities.LTP, error) {
	ret, err := l.source.GetData(context.TODO(), l.pairs[0])
	if err != nil {
		return nil, err
	}

	return ret, nil
}
