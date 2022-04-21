package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/routes"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func main() {
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	routes.SetPrometheus(r) // Set up Prometheus.
	routes.SetRoutes(r)     // Set up all API routes.

	// Run the server
	PORT := zutils.GetEnv("PORT", "4080")
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: r,
	}

	shutdown(func(grace bool) {
		core.CloseIndexes() // close all indexes
		log.Info().Msgf("Index closed")
		if grace {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Fatal().Msgf("Server Shutdown:", err)
			}
		} else {
			server.Close()
		}
	})

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info().Msgf("Server closed under request")
		} else {
			log.Fatal().Msgf("Server closed unexpect")
		}
	}

	log.Info().Msgf("Server shutdown ok")
}

//shutdown support twice signal must exit
func shutdown(stop func(grace bool)) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sig
		go stop(s != syscall.SIGQUIT)
		<-sig
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
}
