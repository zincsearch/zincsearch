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

func TestToString(t *testing.T) {
	type args struct {
		v interface{}
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
				v: "test",
			},
			want: "test",
		},
		{
			name: "float64",
			args: args{
				v: 3.14,
			},
			want: "3.14",
		},
		{
			name: "uint64",
			args: args{
				v: uint64(3),
			},
			want: "3",
		},
		{
			name: "int64",
			args: args{
				v: int64(3),
			},
			want: "3",
		},
		{
			name: "int",
			args: args{
				v: int(3),
			},
			want: "3",
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToString(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				v: "3.14",
			},
			want: 3.14,
		},
		{
			name: "float64",
			args: args{
				v: 3.14,
			},
			want: 3.14,
		},
		{
			name: "uint64",
			args: args{
				v: uint64(3),
			},
			want: float64(3),
		},
		{
			name: "int64",
			args: args{
				v: int64(3),
			},
			want: float64(3),
		},
		{
			name: "int",
			args: args{
				v: int(3),
			},
			want: float64(3),
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: float64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				v: "3",
			},
			want: uint64(3),
		},
		{
			name: "float64",
			args: args{
				v: 3.14,
			},
			want: uint64(3),
		},
		{
			name: "uint64",
			args: args{
				v: uint64(3),
			},
			want: uint64(3),
		},
		{
			name: "int64",
			args: args{
				v: int64(3),
			},
			want: uint64(3),
		},
		{
			name: "int",
			args: args{
				v: int(3),
			},
			want: uint64(3),
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: uint64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint64(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				v: "3",
			},
			want: int(3),
		},
		{
			name: "float64",
			args: args{
				v: 3.14,
			},
			want: int(3),
		},
		{
			name: "uint64",
			args: args{
				v: uint64(3),
			},
			want: int(3),
		},
		{
			name: "int64",
			args: args{
				v: int64(3),
			},
			want: int(3),
		},
		{
			name: "int",
			args: args{
				v: int(3),
			},
			want: int(3),
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: int(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				v: "true",
			},
			want: true,
		},
		{
			name: "float64",
			args: args{
				v: 3.14,
			},
			want: true,
		},
		{
			name: "uint64",
			args: args{
				v: uint64(3),
			},
			want: true,
		},
		{
			name: "int64",
			args: args{
				v: int64(3),
			},
			want: true,
		},
		{
			name: "int",
			args: args{
				v: int(3),
			},
			want: true,
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
