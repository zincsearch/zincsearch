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

import "testing"

func TestStringToInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should convert numeric string to int",
			args: args{s: "1"},
			want: 1,
		},
		{
			name: "should convert numeric string containing whitespace to int",
			args: args{s: "		1 "},
			want: 1,
		},
		{
			name: "should return 0 when given a non numeric string",
			args: args{s: "zero"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToInt(tt.args.s); got != tt.want {
				t.Errorf("StringToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should view '123' as numeric",
			args: args{"123"},
			want: true,
		},
		{
			name: "should view '023' as numeric",
			args: args{"023"},
			want: true,
		},
		{
			name: "should not view '02 3' as numeric",
			args: args{"02 3"},
			want: false,
		},
		{
			name: "should not view 'abc' as numeric",
			args: args{"abc"},
			want: false,
		},
		{
			name: "should view unicode '２' as numeric",
			args: args{"２"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumeric(tt.args.s); got != tt.want {
				t.Errorf("IsNumeric(%v) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}

func TestSliceExists(t *testing.T) {
	type args struct {
		set  []string
		find string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should not find anything in empty set []",
			args: args{
				set:  []string{},
				find: "abc",
			},
			want: false,
		},
		{
			name: "should find abc in [def,abc] set",
			args: args{
				set:  []string{"def", "abc"},
				find: "abc",
			},
			want: true,
		},
		{
			name: "should not find abc in [def,ghi] set",
			args: args{
				set:  []string{"def", "ghi"},
				find: "abc",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceExists(tt.args.set, tt.args.find); got != tt.want {
				t.Errorf("SliceExists(%v, %v) = %v, want %v", tt.args.set, tt.args.find, got, tt.want)
			}
		})
	}
}
