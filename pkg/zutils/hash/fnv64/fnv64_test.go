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

package fnv64

import "testing"

func Test_fnv64a_Sum64(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "test1",
			args: args{
				key: "test1",
			},
			want: 2271358237066212092,
		},
		{
			name: "test2",
			args: args{
				key: "test2",
			},
			want: 2271361535601096725,
		},
	}

	f := fnv64a{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f.Sum64(tt.args.key); got != tt.want {
				t.Errorf("fnv64a.Sum64() = %v, want %v", got, tt.want)
			}
		})
	}
}
