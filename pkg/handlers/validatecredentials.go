package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
)

func ValidateCredentials(c *gin.Context) {
	var user auth.ZincUser
	c.BindJSON(&user)

	validationResult, loggedInUser := auth.VerifyCredentials(user.ID, user.Password)
	c.JSON(200, gin.H{
		"validated": validationResult,
		"user":      loggedInUser,
	})
}
