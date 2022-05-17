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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFirstUser(t *testing.T) {
	type args struct {
		init func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "init first user",
			args: args{
				init: func() {},
			},
			wantErr: false,
		},
		{
			name: "init first user with error",
			args: args{
				init: func() {
					os.Setenv("ZINC_FIRST_ADMIN_USER", "")
					os.Setenv("ZINC_FIRST_ADMIN_PASSWORD", "")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.init()
			err := initFirstUser()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
