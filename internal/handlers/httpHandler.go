package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rmarku/ltp_api/internal/domain"
	"github.com/rmarku/ltp_api/internal/entities"
)

var _ HTTPHandler = new(HTTPHandlerImpl)

func NewHTTPHandler(router *gin.RouterGroup, service domain.LastTradePrice) *HTTPHandlerImpl {
	return &HTTPHandlerImpl{
		router:     router,
		ltpService: service,
	}
}

func (h *HTTPHandlerImpl) Register() {
	h.router.GET("/ltp", h.getLTP)
}

func (h *HTTPHandlerImpl) getLTP(c *gin.Context) {
	var prices []entities.LTP //nolint: prealloc

	var pairs []string

	query := c.Query("pairs")

	if query != "" {
		pairs = strings.Split(query, ",")
	} else {
		pairs = h.ltpService.GetPairs()
	}

	for _, pair := range pairs {
		price, err := h.ltpService.GetLastTradePrices(pair)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Cannot get last trade price",
			})

			slog.Error("Cannot get last trade price", "err", err)

			return
		}

		prices = append(prices, *price)
	}

	c.JSON(http.StatusOK, gin.H{
		"ltp": prices,
	})
}
