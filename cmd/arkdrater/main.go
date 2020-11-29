package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thmhoag/arkdrater/handler"
	"github.com/thmhoag/arkdrater/pkg/arkdrater/dynamic"

	"github.com/thmhoag/arkdrater/pkg/arkdrater/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	log.Println("starting up...")
	log.Println("loading config...")
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalln("unable to load config", err)
	}

	defaultReqHandler := &handler.DefaultHandler{
		RateConverter: dynamic.NewRateConverter(http.DefaultClient),
		Config: cfg,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultReqHandler.HandleRequest)

	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// if your BaseContext is more complex you might want to use this instead of doing it manually
	// httpServer.RegisterOnShutdown(cancel)

	// Run server
	go func() {
		log.Printf("listening on port %v\n", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// it is fine to use Fatal here because it is not main gorutine
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	log.Println("shutting down...")

	go func() {
		<-signalChan
		log.Fatalln("terminating...")
	}()

	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefulCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	}

	cancel()

	defer os.Exit(0)
	return
}
