package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
)

func CreateUpdateUser(c *gin.Context) {
	var user auth.ZincUser
	c.BindJSON(&user)
	if user.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user.id should be not empty"})
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
