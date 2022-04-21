package auth

import (
	"context"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
)

func GetUser(userID string) (ZincUser, bool, error) {
	userExists := false
	var user ZincUser

	query := bluge.NewTermQuery(userID)
	searchRequest := bluge.NewTopNSearch(1, query)
	reader, _ := core.ZINC_SYSTEM_INDEX_LIST["_users"].Writer.Reader()
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("auth.GetUser: error executing search: %s", err.Error())
		return user, userExists, err
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
			log.Printf("auth.GetUser: error accessing stored fields: %s", err.Error())
			return user, userExists, err
		} else {
			return user, userExists, nil
		}
	}

	return user, userExists, nil
}
