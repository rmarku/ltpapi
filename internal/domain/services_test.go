package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rmarku/ltp_api/internal/datasources"
	"github.com/rmarku/ltp_api/internal/entities"
	"github.com/rmarku/ltp_api/internal/keyvalue"
)

func TestGetLastTradePrices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		cacheGetError error
		cacheSetError error
		sourceError   error
		expectedError error
		expected      *entities.LTP
		name          string
		cacheSetValue float64
	}{
		{
			name:          "Cache hit",
			cacheSetValue: 100.0,
			expected: &entities.LTP{
				Pair:   "BTC/USD",
				Amount: 100.0,
			},
		},
		{
			name:          "Cache miss",
			cacheGetError: keyvalue.ErrExpired,
			cacheSetValue: 200.0,
			expected: &entities.LTP{
				Pair:   "BTC/USD",
				Amount: 200.0,
			},
		},
		{
			name:          "Cache error",
			cacheGetError: errors.New("cache error"), //nolint: goerr113
			cacheSetValue: 0,
			expectedError: errors.New("cache error"), //nolint: goerr113
		},
		{
			name:          "Source error",
			cacheGetError: keyvalue.ErrExpired,
			cacheSetError: errors.New("source error"), //nolint: goerr113
			cacheSetValue: 0,
			expected:      nil,
			expectedError: errors.New("source error"), //nolint: goerr113
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSource := new(datasources.MockDataSource)
			mockCache := new(keyvalue.MockFloatCache)

			ltp := NewLastTradePrice(mockSource, mockCache)

			mockCache.On("Get", "BTC/USD").Return(tt.cacheSetValue, tt.cacheGetError)

			if tt.cacheGetError != nil && errors.Is(tt.cacheGetError, keyvalue.ErrExpired) {
				mockSource.On("GetData", mock.Anything, "BTC/USD").
					Return(&entities.LTP{Pair: "BTC/USD", Amount: tt.cacheSetValue}, tt.sourceError)
				mockCache.On("Set", "BTC/USD", mock.AnythingOfType("float64")).Return(tt.cacheSetError)
			}

			result, err := ltp.GetLastTradePrices("BTC/USD")

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedError, err)

			mockCache.AssertExpectations(t)
			mockSource.AssertExpectations(t)
		})
	}
}
