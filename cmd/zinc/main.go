package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	zerolog "github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/routes"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		zerolog.Print("Error loading .env file")
	}

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	routes.SetRoutes(r) // Set up all API routes.

	// Run the server

	PORT := zutils.GetEnv("PORT", "4080")

	r.Run(":" + PORT)
}
