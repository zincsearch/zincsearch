package routes

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
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

	// get global Monitor object
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(5)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	// set middleware for gin
	m.Use(r)
}
