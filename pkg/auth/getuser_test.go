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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    ZincUser
		want1   bool
		wantErr bool
		input   ZincUser
	}{
		{
			name: "get user",
			args: args{
				userID: "testuser",
			},
			want: ZincUser{
				ID:   "testuser",
				Name: "Test User",
				Role: "admin",
			},
			want1: true,
			input: ZincUser{
				ID:       "testuser",
				Name:     "Test User",
				Role:     "admin",
				Password: "testpassword",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(os.Getwd())
			CreateUser(tt.input.ID, tt.input.Name, tt.input.Password, tt.input.Role)
			got, got1, err := GetUser(tt.args.userID)

			assert.Equal(t, err, nil)
			assert.Equal(t, got.ID, tt.want.ID)
			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.Role, tt.want.Role)
			assert.Equal(t, got1, tt.want1)

			// os.RemoveAll("data")
		})
	}
}
