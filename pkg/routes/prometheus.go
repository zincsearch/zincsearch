package routes

import (
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/zinclabs/zinc/pkg/zutils"
)

// SetPrometheus sets up prometheus metrics for gin
func SetPrometheus(r *gin.Engine) {
	if !zutils.GetEnvToBool("ZINC_PROMETHEUS_ENABLE", "false") {
		return
	}

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
}
