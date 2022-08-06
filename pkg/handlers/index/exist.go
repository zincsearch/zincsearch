package index

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Id Exists
// @Summary Checks if the index exists
// @Tags    Index
// @Produce json
// @Param   index  path  string  true  "Index"
// @Success 200 {object} meta.HTTPResponse
// @Failure 404 {object} meta.HTTPResponse
// @Router /api/index/{index} [head]
func Exists(c *gin.Context) {
	indexName := c.Param("target")

	_, exists := core.GetIndex(indexName)
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "ok"})
}

// @Id EsExists
// @Summary Checks if the index exists for compatible ES
// @Tags    Index
// @Produce json
// @Param   index  path  string  true  "Index"
// @Success 200 {object} meta.HTTPResponse
// @Failure 404 {object} meta.HTTPResponse
// @Router /es/{index} [head]
func ESExists() {}
