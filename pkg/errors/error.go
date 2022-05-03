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

import "fmt"

const (
	ErrorTypeParsingException         = "parsing_exception"
	ErrorTypeXContentParseException   = "x_content_parse_exception"
	ErrorTypeIllegalArgumentException = "illegal_argument_exception"
	ErrorTypeNotImplemented           = "not_implemented"
	ErrorTypeRuntimeException         = "runtime_exception"
)

type Error struct {
	Type     string `json:"type"`
	Reason   string `json:"reason"`
	CausedBy error  `json:"caused_by,omitempty"`
}

func New(errType string, errReason string) *Error {
	return &Error{Type: errType, Reason: errReason}
}

func (e *Error) Cause(err error) *Error {
	e.CausedBy = err
	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("error_type: %s, reason: %s", e.Type, e.Reason)
}
