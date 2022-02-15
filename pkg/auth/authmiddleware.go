package auth

import (
	"context"
	"net/http"

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
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

func VerifyCredentials(user, password string) (*SimpleUser, bool) {
	reader, _ := core.ZINC_SYSTEM_INDEX_LIST["_users"].Writer.Reader()
	defer reader.Close()

	termQuery := bluge.NewTermQuery(user).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, termQuery)
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("auth.VerifyCredentials: error executing search: %v", err)
		return nil, false
	}

	sUser := new(SimpleUser)
	storedSalt := ""
	storedPassword := ""

	next, err := dmi.Next()
	if err != nil {
		log.Printf("auth.VerifyCredentials: error accessing search: %v", err)
		return nil, false
	}
	err = next.VisitStoredFields(func(field string, value []byte) bool {
		switch field {
		case "salt":
			storedSalt = string(value)
		case "password":
			storedPassword = string(value)
		case "role":
			sUser.Role = string(value)
		case "name":
			sUser.Name = string(value)
		case "_id":
			sUser.ID = string(value)
		default:
		}

		return true
	})
	if err != nil {
		log.Printf("auth.VerifyCredentials: error accessing stored fields: %v", err)
		return nil, false
	}

	incomingEncryptedPassword := GeneratePassword(password, storedSalt)
	if incomingEncryptedPassword == storedPassword {
		return sUser, true
	}

	return nil, false
}
