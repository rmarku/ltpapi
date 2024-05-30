package domain

import (
	"github.com/stretchr/testify/mock"

	"github.com/rmarku/ltp_api/internal/entities"
)

// Mock for the LastTradePrice service.
type MockLastTradePrice struct {
	mock.Mock
}

var _ LastTradePrice = new(MockLastTradePrice)

func (m *MockLastTradePrice) GetPairs() []string {
	args := m.Called()

	return args.Get(0).([]string)
}

func (m *MockLastTradePrice) GetLastTradePrices(pair string) (*entities.LTP, error) {
	args := m.Called(pair)

	return args.Get(0).(*entities.LTP), args.Error(1)
}

func (m *MockLastTradePrice) UpdatePrices() error {
	args := m.Called()

	return args.Error(0)
}
