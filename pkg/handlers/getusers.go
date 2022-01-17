package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
)

func GetUsers(c *gin.Context) {
	res, err := auth.GetAllUsersWorker()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}
