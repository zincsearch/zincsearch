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

package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Type     string
		Reason   string
		CausedBy error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "error",
			fields: fields{
				Type:   ErrorTypeRuntimeException,
				Reason: "error message",
			},
			want: "type: runtime_exception, reason: error message",
		},
		{
			name: "error caused by",
			fields: fields{
				Type:     ErrorTypeRuntimeException,
				Reason:   "error message",
				CausedBy: errors.New("reason"),
			},
			want: "type: runtime_exception, reason: error message, cause: reason",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Type:     tt.fields.Type,
				Reason:   tt.fields.Reason,
				CausedBy: tt.fields.CausedBy,
			}
			assert.Equal(t, tt.want, e.Error())
		})
	}
}
