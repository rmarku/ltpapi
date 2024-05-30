package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/rmarku/ltp_api/internal/domain"
	"github.com/rmarku/ltp_api/internal/entities"
)

func TestGetLTP(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		mockPrices     map[string]*entities.LTP
		mockErrors     map[string]error
		name           string
		query          string
		expectedBody   string
		mockPairs      []string
		expectedStatus int
	}{
		{
			name:      "No pairs in query, success response",
			query:     "",
			mockPairs: []string{"BTC/USD", "ETH/USD"},
			mockPrices: map[string]*entities.LTP{
				"BTC/USD": {Pair: "BTC/USD", Amount: 30000},
				"ETH/USD": {Pair: "ETH/USD", Amount: 2000},
			},
			mockErrors:     map[string]error{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ltp":[{"pair":"BTC/USD","amount":30000},{"pair":"ETH/USD","amount":2000}]}`,
		},
		{
			name:  "Pairs in query, success response",
			query: "BTC/USD,ETH/USD",
			mockPrices: map[string]*entities.LTP{
				"BTC/USD": {Pair: "BTC/USD", Amount: 30000},
				"ETH/USD": {Pair: "ETH/USD", Amount: 2000}},
			mockErrors:     map[string]error{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ltp":[{"pair":"BTC/USD","amount":30000},{"pair":"ETH/USD","amount":2000}]}`,
		},
		{
			name:  "Error getting prices",
			query: "BTC/USD,ETH/USD",
			mockPrices: map[string]*entities.LTP{
				"BTC/USD": {Pair: "BTC/USD", Amount: 30000},
			},
			mockErrors:     map[string]error{"ETH/USD": domain.ErrPriceNotFound},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := new(domain.MockLastTradePrice)

			if tt.query == "" {
				mockService.On("GetPairs").Return(tt.mockPairs)
			}

			for pair, price := range tt.mockPrices {
				mockService.On("GetLastTradePrices", pair).Return(price, nil)
			}

			for pair, err := range tt.mockErrors {
				mockService.On("GetLastTradePrices", pair).Return(nil, err)
			}

			router := gin.Default()
			handler := NewHTTPHandler(router.Group("/"), mockService)
			handler.Register()

			req := httptest.NewRequest(http.MethodGet, "/ltp?pairs="+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
