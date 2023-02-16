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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zincsearch/test/utils"
)

func TestHandleError(t *testing.T) {
	type args struct {
		err    error
		code   int
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "error",
			args: args{
				err:    errors.New("error message"),
				code:   http.StatusBadRequest,
				result: `{"error":"error message"}`,
			},
		},
		{
			name: "errorx",
			args: args{
				err:    New(ErrorTypeRuntimeException, "error message"),
				code:   http.StatusBadRequest,
				result: `{"error":{"type":"runtime_exception","reason":"error message"}}`,
			},
		},
		{
			name: "errorx with cause",
			args: args{
				err:    New(ErrorTypeRuntimeException, "error message").Cause(errors.New("reason")),
				code:   http.StatusBadRequest,
				result: `{"error":{"type":"runtime_exception","reason":"error message","cause":"reason"}}`,
			},
		},
		{
			name: "nil",
			args: args{
				err:    nil,
				code:   http.StatusOK,
				result: `{"message":"ok"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			HandleError(c, tt.args.err)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Equal(t, tt.args.result, w.Body.String())
		})
	}
}
