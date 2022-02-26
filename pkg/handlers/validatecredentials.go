package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/auth"
)

func ValidateCredentials(c *gin.Context) {
	var user auth.ZincUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loggedInUser, validationResult := auth.VerifyCredentials(user.ID, user.Password)
	c.JSON(http.StatusOK, gin.H{
		"validated": validationResult,
		"user":      loggedInUser,
	})
}
