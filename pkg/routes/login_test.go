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

package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/test/utils"
)

func TestLogin(t *testing.T) {
	type result struct {
		authenticated bool
		errMsg        string
	}
	type args struct {
		code   int
		data   map[string]interface{}
		result result
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should authenticate user",
			args: args{
				data: map[string]interface{}{
					"_id":      "admin",
					"password": "Complexpass#123",
				},
				result: result{
					authenticated: true,
				},
			},
		},
		{
			name: "should not authenticate user with incorrect credentials",
			args: args{
				data: map[string]interface{}{
					"_id":      "user_not_exists",
					"password": "password",
				},
				result: result{
					authenticated: false,
					errMsg:        "Invalid credentials",
				}},
		},
		{
			name: "should not authenticate user with invalid credentials structure",
			args: args{
				data: map[string]interface{}{
					"_id":      1233,
					"password": "password",
				},
				result: result{
					authenticated: false,
					errMsg:        "Invalid credentials structure",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)

			user, err := Login(c)

			if tt.args.result.authenticated {
				assert.NoError(t, err)
				assert.Equal(t, LoginUser{
					ID:   "admin",
					Name: "admin",
					Role: "admin",
				}, user)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.args.result.errMsg)
			}
		})
	}
}
