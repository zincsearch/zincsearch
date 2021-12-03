package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealthz function gets all events
func GetHealthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
