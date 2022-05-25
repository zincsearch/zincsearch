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
	"time"

	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
)

func GetAllUsersWorker() (*meta.SearchResponse, error) {
	users, err := metadata.User.List(0, 0)
	if err != nil {
		return nil, err
	}

	var Hits []meta.Hit
	for _, u := range users {
		hit := meta.Hit{
			Index:     u.Name,
			Type:      u.Name,
			ID:        u.ID,
			Timestamp: u.UpdatedAt,
			Source: SimpleUser{
				ID:        u.ID,
				Name:      u.Name,
				Role:      u.Role,
				Salt:      u.Salt,
				Password:  u.Password,
				CreatedAt: u.CreatedAt,
				Timestamp: u.UpdatedAt,
			},
		}
		Hits = append(Hits, hit)
	}

	resp := &meta.SearchResponse{
		Took: 0,
		Hits: meta.Hits{
			Total: meta.Total{
				Value: len(users),
			},
			MaxScore: 0,
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
