package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ZincAuthMiddleware(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		if _, ok := VerifyCredentials(user, password); ok {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth": "Invalid credentials"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth": "Missing credentials"})
		return
	}
}

func VerifyCredentials(userID, password string) (SimpleUser, bool) {
	user, ok := ZINC_CACHED_USERS[userID]
	if !ok {
		return SimpleUser{}, false
	}

	incomingEncryptedPassword := GeneratePassword(password, user.Salt)
	if incomingEncryptedPassword == user.Password {
		return user, true
	}

	return SimpleUser{}, false
}
