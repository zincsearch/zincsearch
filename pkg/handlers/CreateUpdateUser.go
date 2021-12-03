package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	auth "github.com/prabhatsharma/zinc/pkg/auth"
)

func CreateUpdateUser(c *gin.Context) {
	fmt.Println("CreateUser")

	var user auth.ZincUser
	c.BindJSON(&user)

	newUser, err := auth.CreateUser(user.ID, user.Name, user.Password, user.Role)

	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": newUser,
	})
}
