package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zinclabs/zinc/pkg/auth"
)

func DeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"deleted": auth.DeleteUser(c.Param("userID")),
	})
}
