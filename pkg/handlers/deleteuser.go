package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
)

func DeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"deleted": auth.DeleteUser(c.Param("userID")),
	})
}
