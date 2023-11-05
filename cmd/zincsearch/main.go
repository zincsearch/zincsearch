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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/rs/zerolog/log"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/core"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/metadata"
	"github.com/zincsearch/zincsearch/pkg/routes"
)

// @title           Zinc Search engine API
// @version         1.0.0
// @description     Zinc Search engine API documents https://zincsearch-docs.zinc.dev
// @termsOfService  http://swagger.io/terms/

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @contact.name    Zinc Search
// @contact.url     https://www.zincsearch.com

// @securityDefinitions.basic BasicAuth
// @BasePath        /
// @schemes http https
func main() {
	// Version
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("zinc version %s\n", meta.Version)
		os.Exit(0)
	}
	log.Info().Msgf("Starting Zinc %s", meta.Version)

	// Initialize telemetry
	telemetry()
	// Initialize sentry
	sentries()
	// Continuous profiling
	profiling()

	// HTTP init
	app := gin.New()
	routes.Setup(app)

	// Run the server
	PORT := config.Global.ServerPort
	ADDRESS := config.Global.ServerAddress
	server := &http.Server{
		Addr:    ADDRESS + ":" + PORT,
		Handler: app,
	}

	done := shutdown(func(grace bool, done chan<- struct{}) {
		// close http server
		if grace {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Fatal().Err(err).Msg("Server Shutdown")
			}
		} else {
			server.Close()
		}

		log.Info().Msg("Index closing...")
		// close indexes
		err := core.ZINC_INDEX_LIST.Close()
		log.Info().Err(err).Msgf("Index closed")
		// close metadata
		err = metadata.Close()
		log.Info().Err(err).Msgf("Metadata closed")

		done <- struct{}{}
	})

	err := func() error {
		log.Info().Msg("Listen on " + server.Addr)

		if config.Global.ServerTLSCertificateFile != "" && config.Global.ServerTLSKeyFile != "" {

			// set minimum TLS version (1.2) and intermediate cipher suites
			server.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			}

			certFile := config.Global.ServerTLSCertificateFile
			keyFile := config.Global.ServerTLSKeyFile

			return server.ListenAndServeTLS(certFile, keyFile)

		}

		return server.ListenAndServe()
	}()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Info().Msg("Server closed")
		} else {
			log.Fatal().Err(err).Msg("Server closed unexpect")
		}
	}

	<-done

	log.Info().Msg("Server shutdown ok")
}

func telemetry() {
	core.Telemetry.Instance()
	core.Telemetry.Event("server_start", nil)
	core.Telemetry.Cron()
}

func sentries() {
	if !config.Global.SentryEnable {
		return
	}
	if config.Global.SentryDSN == "" {
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:     config.Global.SentryDSN,
		Release: "zinc@" + meta.Version,
	})
	if err != nil {
		log.Print("sentry.Init: ", err.Error())
	}
}

func profiling() {
	if !config.Global.ProfilerEnable {
		return
	}
	if config.Global.ProfilerServer == "" {
		return
	}

	ProfileID := config.Global.ProfilerFriendlyProfileID
	if ProfileID == "" {
		ProfileID = strings.ToLower(core.Telemetry.GetInstanceID())
	}

	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "zincsearch-" + ProfileID,

		// replace this with the address of pyroscope server
		ServerAddress: config.Global.ProfilerServer,

		// you can disable logging by setting this to nil
		// Logger: pyroscope.StandardLogger,
		Logger: nil,

		// optionally, if authentication is enabled, specify the API key:
		// AuthToken: os.Getenv("PYROSCOPE_AUTH_TOKEN"),
		AuthToken: config.Global.ProfilerAPIKey,

		// by default all profilers are enabled,
		// but you can select the ones you want to use:
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
	if err != nil {
		log.Print("pyroscope.Start: ", err.Error())
	}
}

// shutdown support twice signal must exit
func shutdown(stop func(grace bool, done chan<- struct{})) <-chan struct{} {
	done := make(chan struct{})
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sig
		go stop(s != syscall.SIGQUIT, done)
		<-sig
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
	return done
}
