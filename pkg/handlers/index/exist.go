package index

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
)

// @Summary Checks if the index exists for compatible ES
// @Tags    Index
// @Produce json
// @Param   index body meta.IndexSimple true "Index data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} meta.HTTPResponse
// @Router /es/:target [head]
func Exist(c *gin.Context) {
	indexName := c.Param("target")

	_, exists := core.GetIndex(indexName)
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
