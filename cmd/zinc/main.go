/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/routes"
)

// @title           Zinc Search engine API
// @version         1.0
// @description     Zinc Search engine API
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4080
// @BasePath  /api
func main() {

	if config.Global.SentryEnable {
		/******** initialize sentry **********/
		err := sentry.Init(sentry.ClientOptions{
			Dsn:     "https://15b6d9b8be824b44896f32b0234c32b7@o1218932.ingest.sentry.io/6360942",
			Release: "zinc@" + meta.Version,
		})
		if err != nil {
			log.Print("sentry.Init: ", err.Error())
		}
		/******** sentry initialize complete *******/
	}

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	routes.SetPrometheus(r) // Set up Prometheus.
	routes.SetRoutes(r)     // Set up all API routes.

	// Run the server
	PORT := config.Global.ServerPort
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: r,
	}

	shutdown(func(grace bool) {
		core.ZINC_INDEX_LIST.Close() // close all indexes
		log.Info().Msgf("Index closed")
		if grace {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Fatal().Msgf("Server Shutdown: %s", err.Error())
			}
		} else {
			server.Close()
		}
	})

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info().Msgf("Server closed under request")
		} else {
			log.Fatal().Msgf("Server closed unexpect: %s", err.Error())
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
