package routes

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// SetPrometheus sets up prometheus metrics for gin
func SetPrometheus(r *gin.Engine) {
	enable := false
	if v := os.Getenv("ZINC_PROMETHEUS_ENABLE"); v != "" {
		enable, _ = strconv.ParseBool(v)
	}
	if !enable {
		return
	}

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
}
