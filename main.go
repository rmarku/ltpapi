package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rmarku/ltp_api/internal/datasources"
	"github.com/rmarku/ltp_api/internal/domain"
	"github.com/rmarku/ltp_api/internal/handlers"
	"github.com/rmarku/ltp_api/internal/keyvalue"
)

func main() {
	// Global Context to close application
	ctx, close := context.WithCancel(context.Background())

	// Initialize data source (secondary adapter)
	kraken := datasources.NewKraken("https://api.kraken.com/0/public/Ticker")

	cache := keyvalue.NewInMemory()
	// Initialize services.

	ltpService := domain.NewLastTradePrice(kraken, cache)
	ltpService.UpdatePrices()

	//go StartTicker(ctx, ltpService, time.Minute)

	// Initialize routes for primary adapters
	router := gin.Default()
	api := router.Group("/api/v1")

	// Initialize http handlers (primary adapter)
	apiHandler := handlers.NewHTTPHandler(api, ltpService)

	apiHandler.Register()

	// Start ticker to keep values updated
	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// controlando ctx.Done(). tiempo de espera de 5 segundos.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
}

func StartTicker(ctx context.Context, service domain.LastTradePrice, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutdown Ticker ...")

			return
		case <-ticker.C:
			err := service.UpdatePrices()
			if err != nil {
				slog.Error("error updating prices", "err", err)
			}
		}
	}
}
