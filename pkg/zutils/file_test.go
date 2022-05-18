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

func TestDirSize(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "dir size",
			args: args{
				path: "./",
			},
			want:    float64(0),
			wantErr: false,
		},
		{
			name: "dir not exists",
			args: args{
				path: "testdataNotExists",
			},
			want:    float64(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DirSize(tt.args.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.NoError(t, err)
			assert.Greater(t, got, tt.want)
		})
	}
}

func TestIsExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "is exist",
			args: args{
				path: "./file.go",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "is not exist",
			args: args{
				path: "testdataNotExists",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsExist(tt.args.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
