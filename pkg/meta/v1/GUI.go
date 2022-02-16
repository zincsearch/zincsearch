package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GUI function gets all events
func GUI(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"zinc": "Modern, Simpler, Lighter, Faster Search server. ",
	// })

	c.Status(http.StatusFound)
	c.Writer.Header().Set("Location", "./ui/")
}
