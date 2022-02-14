package auth

import (
	"context"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
)

func GetUser(userId string) (bool, ZincUser, error) {
	userExists := false
	var user ZincUser

	query := bluge.NewTermQuery(userId)
	searchRequest := bluge.NewTopNSearch(1, query)
	usersIndex := core.ZINC_SYSTEM_INDEX_LIST["_users"]
	reader, _ := usersIndex.Writer.Reader()
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
	}

	next, err := dmi.Next()
	for err == nil && next != nil {
		userExists = true
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_id":
				user.ID = string(value)
			case "name":
				user.Name = string(value)
			case "salt":
				user.Salt = string(value)
			case "password":
				user.Password = string(value)
			case "role":
				user.Role = string(value)
			case "created_at":
				user.CreatedAt, _ = bluge.DecodeDateTime(value)
			case "@timestamp":
				user.Timestamp, _ = bluge.DecodeDateTime(value)
			default:
			}

			return true
		})
		if err != nil {
			log.Printf("error accessing stored fields: %v", err)
			return userExists, user, err
		} else {
			return userExists, user, nil
		}

		// next, err = dmi.Next()
	}

	return false, user, nil
}
