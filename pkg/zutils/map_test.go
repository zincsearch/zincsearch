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

package zutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStringFromMap(t *testing.T) {
	type args struct {
		m   interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				m: map[string]interface{}{
					"key": "value",
				},
				key: "key",
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "error type",
			args: args{
				m: map[string]interface{}{
					"key": false,
				},
				key: "key",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "not exist",
			args: args{
				m:   map[string]interface{}{},
				key: "key",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStringFromMap(tt.args.m, tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetBoolFromMap(t *testing.T) {
	type args struct {
		m   interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "true",
			args: args{
				m: map[string]interface{}{
					"key": true,
				},
				key: "key",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "false",
			args: args{
				m: map[string]interface{}{
					"key": false,
				},
				key: "key",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error type",
			args: args{
				m: map[string]interface{}{
					"key": "value",
				},
				key: "key",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "not exist",
			args: args{
				m:   map[string]interface{}{},
				key: "key",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBoolFromMap(tt.args.m, tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetFloatFromMap(t *testing.T) {
	type args struct {
		m   interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "float",
			args: args{
				m: map[string]interface{}{
					"key": 3.14,
				},
				key: "key",
			},
			want:    3.14,
			wantErr: false,
		},
		{
			name: "error type",
			args: args{
				m: map[string]interface{}{
					"key": "value",
				},
				key: "key",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "not exist",
			args: args{
				m:   map[string]interface{}{},
				key: "key",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFloatFromMap(tt.args.m, tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetStringSliceFromMap(t *testing.T) {
	type args struct {
		m   interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "slice",
			args: args{
				m: map[string]interface{}{
					"key": []string{"value1", "value2"},
				},
				key: "key",
			},
			want:    []string{"value1", "value2"},
			wantErr: false,
		},
		{
			name: "interface",
			args: args{
				m: map[string]interface{}{
					"key": []interface{}{"value1", "value2"},
				},
				key: "key",
			},
			want:    []string{"value1", "value2"},
			wantErr: false,
		},
		{
			name: "interface error type",
			args: args{
				m: map[string]interface{}{
					"key": []interface{}{"value1", 3.14},
				},
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error type",
			args: args{
				m: map[string]interface{}{
					"key": "value",
				},
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not exist",
			args: args{
				m:   map[string]interface{}{},
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStringSliceFromMap(tt.args.m, tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetMapFromMap(t *testing.T) {
	type args struct {
		m   interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "map",
			args: args{
				m: map[string]interface{}{
					"key": map[string]interface{}{
						"subkey": "value",
					},
				},
				key: "key",
			},
			want:    map[string]interface{}{"subkey": "value"},
			wantErr: false,
		},
		{
			name: "error type",
			args: args{
				m: map[string]interface{}{
					"key": "value",
				},
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not exist",
			args: args{
				m:   map[string]interface{}{},
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMapFromMap(tt.args.m, tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetAnyFromMap(t *testing.T) {
	type args struct {
		m   interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "exist",
			args: args{
				m: map[string]interface{}{
					"key": "value",
				},
				key: "key",
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				m:   map[string]interface{}{},
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil",
			args: args{
				m:   nil,
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not map",
			args: args{
				m:   3.14,
				key: "key",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAnyFromMap(tt.args.m, tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
