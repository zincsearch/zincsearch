package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc"
	"github.com/prabhatsharma/zinc/pkg/auth"
	"github.com/prabhatsharma/zinc/pkg/handlers"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/rs/zerolog/log"
)

// SetRoutes sets up all gi HTTP API endpoints that can be called by front end
func SetRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// meta service - healthz
	r.GET("/healthz", v1.GetHealthz)
	r.GET("/", v1.GUI)
	r.GET("/version", v1.GetVersion)

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

	// elastic compatible APIs
	// Deprecated - /es/*  will be removed from zinc in future releases and replaced with /api/*
	// Document APIs - https://www.elastic.co/guide/en/elasticsearch/reference/current/docs.html
	// Single document APIs

	// Index - https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-index_.html

	r.PUT("/es/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)

	r.DELETE("/es/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.DeleteDocument)

	r.POST("/es/:target/_doc", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.PUT("/es/:target/_create/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/es/:target/_create/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)

	// Update
	r.POST("/es/:target/_update/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)

	// Bulk update/insert

	r.POST("/es/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)
	r.POST("/es/:target/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)

}
