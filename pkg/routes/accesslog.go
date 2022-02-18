package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func AccessLog(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		timeStart := time.Now()
		c.Writer.Header().Set("Zinc", v1.Version)

		c.Next()

		took := time.Since(timeStart) / time.Millisecond
		log.Info().
			Str("method", c.Request.Method).
			Int("code", c.Writer.Status()).
			Int("took", int(took)).
			Msg(c.Request.RequestURI)
	})
}
