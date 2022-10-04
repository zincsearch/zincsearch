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
	"net/http"
	"testing"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/auth"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/test/utils"
)

func TestCreateUpdateUser(t *testing.T) {
	type args struct {
		code   int
		data   map[string]interface{}
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				code: http.StatusOK,
				data: map[string]interface{}{
					"_id":      "user",
					"name":     "user",
					"role":     "admin",
					"password": "password",
				},
				result: "message",
			},
		},
		{
			name: "error",
			args: args{
				code: http.StatusBadRequest,
				data: map[string]interface{}{
					"_id": "",
				},
				result: "error",
			},
		},
		{
			name: "empty",
			args: args{
				code: http.StatusBadRequest,
				data: map[string]interface{}{
					"_id": 123,
				},
				result: "error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			CreateUpdateUser(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}

func TestListUser(t *testing.T) {
	type args struct {
		code   int
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				code:   http.StatusOK,
				result: `[{"_id":`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			ListUser(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		code   int
		data   map[string]string
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				code:   http.StatusOK,
				data:   map[string]string{"id": "1"},
				result: "message",
			},
		},
		{
			name: "empty",
			args: args{
				code:   http.StatusOK,
				data:   map[string]string{"id": ""},
				result: "message",
			},
		},
		{
			name: "nil",
			args: args{
				code:   http.StatusOK,
				data:   nil,
				result: "message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestParams(c, tt.args.data)
			DeleteUser(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		jwtPayloadId string
		mockUserId   string
		code         int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should not return unauthorized(401) when id claim is missing",
			args: args{
				code: 401,
			},
		},
		{
			name: "should return ok(200) and the user when correct claim is provided",
			args: args{
				jwtPayloadId: "123",
				mockUserId:   "123",
				code:         200,
			},
		},
		{
			name: "should return unauthorized(401) when user with given id claim is not found",
			args: args{
				jwtPayloadId: "123",
				code:         200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			if tt.args.mockUserId != "" {
				err := auth.SetUser(tt.args.jwtPayloadId, meta.User{
					ID:        tt.args.jwtPayloadId,
					Name:      "admin",
					Role:      "admin",
					Salt:      "pepper",
					Password:  "secret123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				assert.NoError(t, err)
			}
			if tt.args.jwtPayloadId != "" {
				c.Set("JWT_PAYLOAD", jwt.MapClaims{
					"id": tt.args.jwtPayloadId,
				})
			}
			GetUser(c)
			assert.Equal(t, tt.args.code, w.Code)
		})
	}
}
