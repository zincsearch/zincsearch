package auth

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/core"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/rs/zerolog/log"
)

func GetAllUsersWorker() (v1.SearchResponse, error) {
	usersIndex := core.ZINC_SYSTEM_INDEX_LIST["_users"]
	var Hits []v1.Hit

	query := bluge.NewMatchAllQuery()

	searchRequest := bluge.NewTopNSearch(1000, query).WithStandardAggregations()

	reader, _ := usersIndex.Writer.Reader()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
	}

	next, err := dmi.Next()
	for err == nil && next != nil {
		var user SimpleUser
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				user.ID = string(value)
				return true
			} else if field == "name" {
				user.Name = string(value)
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
			log.Printf("error accessing stored fields: %v", err)
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

	resp := v1.SearchResponse{
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
	CreatedAt time.Time `json:"created_at"`
	Timestamp time.Time `json:"@timestamp"`
}
