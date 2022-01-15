package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
)

func DeleteUser(c *gin.Context) {

	userID := c.Param("userID")

	c.JSON(200, gin.H{
		"deleted": auth.DeleteUser(userID),
	})
}
