package server

import (
	"bspfp/bsrevproxy/internal/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer() {
	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", config.Value.Host, config.Value.Port),
		Handler: http.HandlerFunc(requestHandler),
	}

	log.Println("server listening on", server.Addr)

	go func() {
		var err error
		if config.Value.CertFile == "" && config.Value.KeyFile == "" {
			err = server.ListenAndServe()
		} else {
			err = server.ListenAndServeTLS(config.Value.CertFile, config.Value.KeyFile)
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s:%d: %v\n", config.Value.Host, config.Value.Port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %+v\n", err)
	}

	log.Println("server exiting")
}
