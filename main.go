package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/rmarku/ltp_api/internal/config"
	"github.com/rmarku/ltp_api/internal/datasources"
	"github.com/rmarku/ltp_api/internal/domain"
	"github.com/rmarku/ltp_api/internal/handlers"
	"github.com/rmarku/ltp_api/internal/keyvalue"
)

const (
	ExitStatusCode1   = 1
	ExitStatusCode2   = 2
	TimeOutExit       = 10
	TimeOutServer     = 5
	SignalChannelSize = 2
	ReadHeaderTimeout = 10
)

func main() {
	// Global Context to cancel application
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize configuration
	config.InitConfig()

	// Initialize data source (secondary adapter)
	kraken := datasources.NewKraken("https://api.kraken.com/0/public/Ticker")
	cache := keyvalue.NewInMemory()

	// Initialize services.
	ltpService := domain.NewLastTradePrice(kraken, cache)

	if viper.GetBool("ticker.enabled") { // Start ticker if enabled
		StartTicker(ctx, ltpService, viper.GetDuration("ticker.timeout"))
	}

	// Initialize routes for primary adapters
	router := gin.Default()
	api := router.Group("/api/v1")

	// Initialize http handlers (primary adapter)
	apiHandler := handlers.NewHTTPHandler(api, ltpService)

	apiHandler.Register()

	srv := startServer(router, viper.GetString("server.port"))

	go gracefulShutDown(cancel)

	<-ctx.Done()

	stopServer(srv, TimeOutServer*time.Second)
}

func gracefulShutDown(cancel context.CancelFunc) {
	// Wait for SIGTERM and SIGINT system signals
	signalReceived := make(chan os.Signal, SignalChannelSize)
	signal.Notify(signalReceived, syscall.SIGTERM, syscall.SIGINT)
	slog.Debug("GracefulStopService waiting for system signals")

	// Wait for signal to occur once
	sig := <-signalReceived

	slog.Info("caught sig. Started timeout for exit", "signal", sig.String())

	// Inform application that we are closing
	cancel()

	// Wait prudent time for exit, if not send code error
	go func() {
		time.Sleep(TimeOutExit * time.Second)
		slog.Error("time out in the shutdown procedure")
		os.Exit(ExitStatusCode1)
	}()

	// Keep waiting second signal for forced shutdown
	<-signalReceived

	slog.Warn("two consecutive signals received, exiting now")

	os.Exit(ExitStatusCode2)
}

func startServer(router *gin.Engine, port string) *http.Server {
	// Add health check route
	router.GET("/_healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:              port,
		Handler:           router,
		ReadHeaderTimeout: ReadHeaderTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
		}
	}()

	return srv
}

func stopServer(srv *http.Server, timeout time.Duration) {
	slog.Info("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown", "err", err)
	}
}

func StartTicker(ctx context.Context, service domain.LastTradePrice, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutdown Ticker")

			return
		case <-ticker.C:
			err := service.UpdatePrices()
			if err != nil {
				slog.Error("error updating prices", "err", err)
			}
		}
	}
}
