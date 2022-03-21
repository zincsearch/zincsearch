package routes

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HTTPCacheForUI(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			if strings.Contains(c.Request.RequestURI, "/ui/assets/") {
				c.Writer.Header().Set("cache-control", "public, max-age=2592000")
				c.Writer.Header().Set("expires", time.Now().Add(time.Hour*24*30).Format(time.RFC1123))
			}
		}

		c.Next()
	})
}
