/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package auth

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

func GetAllUsersWorker() (*meta.SearchResponse, error) {
	usersIndex := core.ZINC_SYSTEM_INDEX_LIST["_users"]

	query := bluge.NewMatchAllQuery()
	searchRequest := bluge.NewTopNSearch(1000, query).WithStandardAggregations()
	reader, _ := usersIndex.Writer.Reader()
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %s", err.Error())
	}

	var Hits []meta.Hit
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

		hit := meta.Hit{
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

	resp := &meta.SearchResponse{
		Took: int(dmi.Aggregations().Duration().Milliseconds()),
		Hits: meta.Hits{
			Total: meta.Total{
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
