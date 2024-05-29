package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/rmarku/ltp_api/internal/domain"
)

type HTTPHandler interface {
	Register()
}

type HTTPHandlerImpl struct {
	router     *gin.RouterGroup
	ltpService domain.LastTradePrice
}
