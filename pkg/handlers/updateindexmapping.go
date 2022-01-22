package handlers

import "github.com/gin-gonic/gin"

func UpdateIndexMapping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
