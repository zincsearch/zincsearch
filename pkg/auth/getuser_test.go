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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zincsearch/zincsearch/pkg/meta"
)

func TestGetUser(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.User
		want1   bool
		wantErr bool
		input   *meta.User
	}{
		{
			name: "get user",
			args: args{
				id: "testuser",
			},
			want: &meta.User{
				ID:   "testuser",
				Name: "Test User",
				Role: "admin",
			},
			want1: true,
			input: &meta.User{
				ID:       "testuser",
				Name:     "Test User",
				Role:     "admin",
				Password: "testpassword",
			},
		},
		{
			name: "get user not exists",
			args: args{
				id: "testuserNotExists",
			},
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != nil {
				got, err := CreateUser(tt.input.ID, tt.input.Name, tt.input.Password, tt.input.Role)
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
			got, got1, err := GetUser(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Role, got.Role)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
