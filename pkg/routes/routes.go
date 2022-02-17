package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc"
	"github.com/prabhatsharma/zinc/pkg/auth"
	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/handlers"
	handlersV2 "github.com/prabhatsharma/zinc/pkg/handlers/v2"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

// SetRoutes sets up all gin HTTP API endpoints that can be called by front end
func SetRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		log.Info().Str("method", c.Request.Method).Msg(c.Request.RequestURI)
		c.Writer.Header().Set("zinc", v1.Version)
		c.Next()
	})

	r.GET("/", v1.GUI)
	r.GET("/version", v1.GetVersion)
	// meta service - healthz
	r.GET("/healthz", v1.GetHealthz)

	front, err := zinc.GetFrontendAssets()
	if err != nil {
		log.Err(err)
	}

	r.StaticFS("/ui", http.FS(front))

	r.POST("/api/login", handlers.ValidateCredentials)

	r.PUT("/api/user", auth.ZincAuthMiddleware, handlers.CreateUpdateUser)
	r.DELETE("/api/user/:userID", auth.ZincAuthMiddleware, handlers.DeleteUser)
	r.GET("/api/users", auth.ZincAuthMiddleware, handlers.GetUsers)

	r.PUT("/api/index", auth.ZincAuthMiddleware, handlers.CreateIndex)
	r.GET("/api/index", auth.ZincAuthMiddleware, handlers.ListIndexes)
	r.DELETE("/api/index/:indexName", auth.ZincAuthMiddleware, handlers.DeleteIndex)

	// Bulk update/insert
	r.POST("/api/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)
	r.POST("/api/:target/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)

	// Document CRUD APIs. Update is same as create.
	r.PUT("/api/:target/document", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/api/:target/_doc", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.PUT("/api/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/api/:target/_search", auth.ZincAuthMiddleware, handlers.SearchIndex)
	r.DELETE("/api/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.DeleteDocument)
	r.GET("/api/:target/_mapping", auth.ZincAuthMiddleware, handlers.GetIndexMapping)
	r.PUT("/api/:target/_mapping", auth.ZincAuthMiddleware, handlers.UpdateIndexMapping)

	// elastic filebeat
	r.GET("/es/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":         "NA",
			"cluster_name": "NA",
			"cluster_uuid": "NA",
			"version": gin.H{
				"number":                              "0.1.1-zinc",
				"build_flavor":                        "default",
				"build_type":                          "NA",
				"build_hash":                          "NA",
				"build_date":                          "2021-12-12T20:18:09.722761972Z",
				"build_snapshot":                      false,
				"lucene_version":                      "NA",
				"minimum_wire_compatibility_version":  "NA",
				"minimum_index_compatibility_version": "NA",
			},
			"tagline": "You Know, for Search",
		})
	})
	r.GET("/es/_license", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"license": gin.H{
				"status": "active",
			},
		})
	})
	r.GET("/es/_xpack", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"build":    gin.H{},
			"features": gin.H{},
			"license":  gin.H{"status": "active"},
		})
	})

	// elastic compatible APIs
	r.POST("/es/:target/_search", auth.ZincAuthMiddleware, handlersV2.SearchIndex)

	r.GET("/es/_index_template", auth.ZincAuthMiddleware, handlersV2.ListIndexTemplate)
	r.PUT("/es/_index_template/:target", auth.ZincAuthMiddleware, handlersV2.UpdateIndexTemplate)
	r.GET("/es/_index_template/:target", auth.ZincAuthMiddleware, handlersV2.GetIndexTemplate)
	r.HEAD("/es/_index_template/:target", auth.ZincAuthMiddleware, handlersV2.GetIndexTemplate)
	r.DELETE("/es/_index_template/:target", auth.ZincAuthMiddleware, handlersV2.DeleteIndexTemplate)

	r.GET("/es/:target/_mapping", auth.ZincAuthMiddleware, handlers.GetIndexMapping)
	r.PUT("/es/:target/_mapping", auth.ZincAuthMiddleware, handlers.UpdateIndexMapping)

	r.POST("/es/:target/_doc", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.PUT("/es/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.PUT("/es/:target/_create/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/es/:target/_create/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/es/:target/_update/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.DELETE("/es/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.DeleteDocument)

	// Bulk update/insert
	r.POST("/es/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)
	r.POST("/es/:target/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)

	core.TelemetryInstance()
	event_data := make(map[string]interface{})
	core.TelemetryEvent("server_start", event_data)
	core.TelemetryCron()
}
