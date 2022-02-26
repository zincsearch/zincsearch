package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
)

func CreateUpdateUser(c *gin.Context) {
	var user auth.ZincUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user.id should be not empty"})
		return
	}

	newUser, err := auth.CreateUser(user.ID, user.Name, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": newUser,
	})
}
