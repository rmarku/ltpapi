package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rmarku/ltp_api/internal/domain"
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
	price, err := h.ltpService.GetLastTradePrices("BTC/USD")
	if err != nil {
		slog.Error("Cannot get last trade price")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Cannot get last trade price",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"price":  price,
		"status": "ok",
	})
}
