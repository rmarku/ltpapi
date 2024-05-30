package domain

import (
	"context"
	"errors"

	"github.com/spf13/viper"

	"github.com/rmarku/ltp_api/internal/datasources"
	"github.com/rmarku/ltp_api/internal/entities"
	"github.com/rmarku/ltp_api/internal/keyvalue"
)

type LastTradePriceImpl struct {
	source datasources.DataSource
	cache  keyvalue.FloatCache
	pairs  []string
}

var _ LastTradePrice = new(LastTradePriceImpl)

func NewLastTradePrice(source datasources.DataSource, cache keyvalue.FloatCache) *LastTradePriceImpl {
	pairs := viper.GetStringSlice("available_pairs")

	ltp := &LastTradePriceImpl{
		pairs:  pairs,
		source: source,
		cache:  cache,
	}

	ltp.UpdatePrices()

	return ltp
}

func (l *LastTradePriceImpl) GetPairs() []string {
	return l.pairs
}

func (l *LastTradePriceImpl) GetAllLastTradePrices() ([]*entities.LTP, error) {
	result := make([]*entities.LTP, 0, len(l.pairs))

	for _, pair := range l.pairs {
		ret, err := l.GetLastTradePrices(pair)
		if err != nil {
			return nil, err
		}

		result = append(result, ret)
	}

	return result, nil
}

func (l *LastTradePriceImpl) GetLastTradePrices(pair string) (*entities.LTP, error) {
	amount, err := l.cache.Get(pair)
	if err != nil {
		if !errors.Is(err, keyvalue.ErrExpired) {
			return nil, err
		}

		ret, err := l.source.GetData(context.TODO(), pair)
		if err != nil {
			return nil, err
		}

		err = l.cache.Set(pair, ret.Amount)
		if err != nil {
			return nil, err
		}

		amount = ret.Amount
	}

	return &entities.LTP{
		Pair:   pair,
		Amount: amount,
	}, nil
}

func (l *LastTradePriceImpl) UpdatePrices() error {
	for _, pair := range l.pairs {
		ret, err := l.source.GetData(context.TODO(), pair)
		if err != nil {
			return err
		}

		err = l.cache.Set(pair, ret.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}
