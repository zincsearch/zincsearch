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

package etcd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/config"
)

func TestMain(m *testing.M) {
	if config.Global.Etcd.Endpoints == nil {
		os.Exit(0)
		return
	}
	os.Exit(m.Run())
}

func Test_etcdStorage_List(t *testing.T) {
	type args struct {
		prefix string
		in1    int
		in2    int
	}
	tests := []struct {
		name    string
		args    args
		wantNum int
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				prefix: "/test/",
			},
			wantNum: 1,
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				prefix: "/notexist/",
			},
			wantNum: 0,
			wantErr: false,
		},
	}

	store := New("/zincsearch/test")
	defer store.Close()
	t.Run("prepare", func(t *testing.T) {
		err := store.Set("/test/foo", []byte("bar"))
		assert.NoError(t, err)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.List(tt.args.prefix, tt.args.in1, tt.args.in2)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.GreaterOrEqual(t, len(got), tt.wantNum)
		})
	}
}

func Test_etcdStorage_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				key: "/test/foo",
			},
			want:    []byte("bar"),
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				key: "/test/notexist",
			},
			want:    nil,
			wantErr: true,
		},
	}

	store := New("/zincsearch/test")
	defer store.Close()
	t.Run("prepare", func(t *testing.T) {
		err := store.Set("/test/foo", []byte("bar"))
		assert.NoError(t, err)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Get(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_etcdStorage_Set(t *testing.T) {
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				key:   "/test/foo",
				value: []byte("bar"),
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				key:   "",
				value: nil,
			},
			wantErr: true,
		},
	}

	store := New("/zincsearch/test")
	defer store.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := store.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("etcdStorage.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdStorage_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				key: "/test/foo",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				key: "",
			},
			wantErr: true,
		},
	}

	store := New("/zincsearch/test")
	defer store.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := store.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("etcdStorage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
