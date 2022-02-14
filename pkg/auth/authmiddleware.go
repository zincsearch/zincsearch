package auth

import (
	"context"

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
)

func ZincAuthMiddleware(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		result, _ := VerifyCredentials(user, password)
		if result {
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{
				"auth": "Invalid credentials",
			})
			return
		}
	} else {
		c.AbortWithStatusJSON(401, gin.H{
			"auth": "Missing credentials",
		})
		return
	}
}

func VerifyCredentials(user, password string) (bool, SimpleUser) {
	var sUser SimpleUser
	reader, _ := core.ZINC_SYSTEM_INDEX_LIST["_users"].Writer.Reader()
	termQuery := bluge.NewTermQuery(user).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1000, termQuery)
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
	}

	storedSalt := ""
	storedPassword := ""

	next, err := dmi.Next()
	for err == nil && next != nil {
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
			log.Printf("error accessing stored fields: %v", err)
		}

		incomingEncryptedPassword := GeneratePassword(password, storedSalt)
		if incomingEncryptedPassword == storedPassword {
			return true, sUser
		}

		next, err = dmi.Next()
	}

	reader.Close()

	return false, sUser
}
