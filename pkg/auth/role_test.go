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
	"github.com/zinclabs/zinc/pkg/meta"
)

func TestCreateRole(t *testing.T) {
	type args struct {
		id         string
		name       string
		permission []string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.Role
		wantErr bool
	}{
		{
			name: "create role",
			args: args{
				id:         "testrole",
				name:       "Test Role",
				permission: []string{"test"},
			},
			want: &meta.Role{
				ID:         "testrole",
				Name:       "Test Role",
				Permission: []string{"test"},
			},
			wantErr: false,
		},
		{
			name: "update exists role",
			args: args{
				id:         "testrole",
				name:       "Test Role Updated",
				permission: []string{"test", "t2"},
			},
			want: &meta.Role{
				ID:         "testrole",
				Name:       "Test Role Updated",
				Permission: []string{"test", "t2"},
			},
			wantErr: false,
		},
		{
			name: "create role with empty id",
			args: args{
				id: "",
			},
			want: &meta.Role{
				ID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateRole(tt.args.id, tt.args.name, tt.args.permission)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Permission, got.Permission)
		})
	}
}

func TestGetRole(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.Role
		want1   bool
		wantErr bool
		input   *meta.Role
	}{
		{
			name: "get role",
			args: args{
				id: "testrole",
			},
			want: &meta.Role{
				ID:         "testrole",
				Name:       "Test Role",
				Permission: []string{"test"},
			},
			want1: true,
			input: &meta.Role{
				ID:         "testrole",
				Name:       "Test Role",
				Permission: []string{"test"},
			},
		},
		{
			name: "get role not exists",
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
				got, err := CreateRole(tt.input.ID, tt.input.Name, tt.input.Permission)
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
			got, got1, err := GetRole(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Permission, got.Permission)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestDeleteRole(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete role testrole",
			args: args{
				id: "testrole",
			},
			wantErr: false,
		},
		{
			name: "delete role not exists role",
			args: args{
				id: "testuserNotExists",
			},
			wantErr: false,
		},
		{
			name: "delete role empty",
			args: args{
				id: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteRole(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
