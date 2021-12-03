package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	zerolog "github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/routes"
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

	// startup.LoadIndexes() // Load bluge indexes from disk.

	// Run the server

	if os.Getenv("PORT") != "" {
		r.Run(":" + os.Getenv("PORT"))
	} else {
		r.Run(":" + "4080")
	}
}
