package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/auth"
)

func GetUsers(c *gin.Context) {
	res, err := auth.GetAllUsersWorker()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
