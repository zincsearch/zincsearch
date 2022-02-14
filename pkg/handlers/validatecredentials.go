package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/auth"
)

func ValidateCredentials(c *gin.Context) {
	var user auth.ZincUser
	c.BindJSON(&user)

	loggedInUser, validationResult := auth.VerifyCredentials(user.ID, user.Password)
	c.JSON(http.StatusOK, gin.H{
		"validated": validationResult,
		"user":      loggedInUser,
	})
}
