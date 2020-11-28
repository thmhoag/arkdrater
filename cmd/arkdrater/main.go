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
	"text/template"
	"time"

	"github.com/thmhoag/arkdrater/pkg/arkdrater/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	log.Println("starting up...")
	log.Println("loading config...")
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalln("unable to load config", err)
	}

	log.Println("confg loaded")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		template, err := template.New("handlerResponse").Parse(handlerTemplate)
		if err != nil {
			log.Printf("error while marshling object. %s \n", err.Error())
			return
		}

		if err := template.Execute(w, cfg.DynamicConfig); err != nil {
			log.Printf("template.Execute error %s \n", err.Error())
			return
		}
	})

	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// if your BaseContext is more complex you might want to use this instead of doing it manually
	// httpServer.RegisterOnShutdown(cancel)

	// Run server
	go func() {
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
		log.Fatalln("Terminating...")
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

const handlerTemplate = `
{{- range $multiplier, $value := .Multipliers }}
{{ $multiplier }}={{ printf "%.1f" $value }}
{{- end }}
`