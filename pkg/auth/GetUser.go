package auth

import (
	"context"

	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/rs/zerolog/log"
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
		log.Print("error executing search: %v", err)
	}

	next, err := dmi.Next()
	for err == nil && next != nil {
		userExists = true

		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				user.ID = string(value)
				return true
			} else if field == "name" {
				user.Name = string(value)
				return true
			} else if field == "salt" {
				user.Salt = string(value)
				return true
			} else if field == "password" {
				user.Password = string(value)
				return true
			} else if field == "role" {
				user.Role = string(value)
				return true
			} else if field == "created_at" {
				user.CreatedAt, _ = bluge.DecodeDateTime(value)
				return true
			} else if field == "@timestamp" {
				user.Timestamp, _ = bluge.DecodeDateTime(value)
				return true
			}
			return true
		})
		if err != nil {
			log.Print("error accessing stored fields: %v", err)
			return userExists, user, err
		} else {
			return userExists, user, nil
		}

		// next, err = dmi.Next()

	}

	return false, user, nil

}
