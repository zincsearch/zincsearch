package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Analyze(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
