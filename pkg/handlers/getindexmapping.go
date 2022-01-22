package handlers

import (
	"github.com/gin-gonic/gin"
)

func GetIndexMapping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
