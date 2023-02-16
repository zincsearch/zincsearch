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

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zincsearch/test/utils"
)

func TestLogin(t *testing.T) {
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
					"_id":      "admin",
					"password": "Complexpass#123",
				},
				result: "validated\":true",
			},
		},
		{
			name: "error",
			args: args{
				code: http.StatusOK,
				data: map[string]interface{}{
					"_id":      "user_not_exists",
					"password": "password",
				},
				result: "validated\":false",
			},
		},
		{
			name: "error",
			args: args{
				code: http.StatusBadRequest,
				data: map[string]interface{}{
					"_id":      1233,
					"password": "password",
				},
				result: "error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			Login(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}
