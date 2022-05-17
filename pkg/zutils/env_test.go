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
)

func TestGetEnv(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "env",
			args: args{
				key:      "ZINC_FIRST_ADMIN_USER",
				fallback: "value",
			},
			want: "admin",
		},
		{
			name: "not exists",
			args: args{
				key:      "notExists",
				fallback: "value",
			},
			want: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvToLower(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "env",
			args: args{
				key:      "ZINC_FIRST_ADMIN_USER",
				fallback: "VALUE",
			},
			want: "admin",
		},
		{
			name: "not exists",
			args: args{
				key:      "notExists",
				fallback: "VALUE",
			},
			want: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvToLower(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvToUpper(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "env",
			args: args{
				key:      "ZINC_FIRST_ADMIN_USER",
				fallback: "value",
			},
			want: "ADMIN",
		},
		{
			name: "not exists",
			args: args{
				key:      "notExists",
				fallback: "value",
			},
			want: "VALUE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvToUpper(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvToBool(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "env",
			args: args{
				key:      "ZINC_PROMETHEUS_ENABLE",
				fallback: "true",
			},
			want: true,
		},
		{
			name: "not exists",
			args: args{
				key:      "notExists",
				fallback: "true",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvToBool(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
