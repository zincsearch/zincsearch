package main

import (
	"github.com/gin-gonic/gin"

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
	r.Run(":" + PORT)
}
