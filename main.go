package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func rootHandler(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.String(http.StatusOK, "Welcome Gin Server")
}
func main() {


	router := gin.Default()
	router.GET("/", rootHandler)

	

	price, err := getTicket("BTC/USD")
	if err != nil {
		panic("fail")
	}
	log.Printf("%v", *price)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
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
