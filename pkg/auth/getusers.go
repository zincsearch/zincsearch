package auth

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func GetAllUsersWorker() (*v1.SearchResponse, error) {
	usersIndex := core.ZINC_SYSTEM_INDEX_LIST["_users"]

	query := bluge.NewMatchAllQuery()
	searchRequest := bluge.NewTopNSearch(1000, query).WithStandardAggregations()
	reader, _ := usersIndex.Writer.Reader()
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %s", err.Error())
	}

	var Hits []v1.Hit
	next, err := dmi.Next()
	for err == nil && next != nil {
		var user SimpleUser
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_id":
				user.ID = string(value)
			case "name":
				user.Name = string(value)
			case "role":
				user.Role = string(value)
			case "salt":
				user.Salt = string(value)
			case "password":
				user.Password = string(value)
			case "created_at":
				user.CreatedAt, _ = bluge.DecodeDateTime(value)
			case "@timestamp":
				user.Timestamp, _ = bluge.DecodeDateTime(value)
			default:
			}

			return true
		})
		if err != nil {
			log.Printf("error accessing stored fields: %s", err.Error())
		}

		hit := v1.Hit{
			Index:     usersIndex.Name,
			Type:      usersIndex.Name,
			ID:        user.ID,
			Score:     next.Score,
			Timestamp: user.Timestamp,
			Source:    user,
		}
		Hits = append(Hits, hit)

		next, err = dmi.Next()
	}

	resp := &v1.SearchResponse{
		Took: int(dmi.Aggregations().Duration().Milliseconds()),
		Hits: v1.Hits{
			Total: v1.Total{
				Value: int(dmi.Aggregations().Count()),
			},
			MaxScore: dmi.Aggregations().Metric("max_score"),
			Hits:     Hits,
		},
	}

	return resp, nil
}

type SimpleUser struct {
	ID        string    `json:"_id"` // this will be email
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Salt      string    `json:"-"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Timestamp time.Time `json:"@timestamp"`
}
