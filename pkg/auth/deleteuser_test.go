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

import "testing"

func TestDeleteUser(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "delete user testuser",
			args: args{
				userID: "testuser",
			},
			want: true,
		},
		{
			name: "delete user not exists user",
			args: args{
				userID: "testuserNotExists",
			},
			want: true,
		},
		{
			name: "delete user empty",
			args: args{
				userID: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteUser(tt.args.userID); got != tt.want {
				t.Errorf("DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
