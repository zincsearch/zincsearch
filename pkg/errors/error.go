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
	"fmt"
	"strings"
)

const (
	ErrorTypeParsingException         = "parsing_exception"
	ErrorTypeXContentParseException   = "x_content_parse_exception"
	ErrorTypeIllegalArgumentException = "illegal_argument_exception"
	ErrorTypeRuntimeException         = "runtime_exception"
	ErrorTypeNotImplemented           = "not_implemented"
	ErrorTypeInvalidArgument          = "invalid_argument"
)

var ErrorIDNotFound = errors.New("id not found")

type Error struct {
	Type     string `json:"type"`
	Reason   string `json:"reason"`
	CausedBy error  `json:"caused_by,omitempty"`
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func New(errType string, errReason string) *Error {
	return &Error{Type: errType, Reason: errReason}
}

func (e *Error) Cause(err error) *Error {
	e.CausedBy = err
	return e
}

func (e *Error) MarshalJSON() ([]byte, error) {
	reason := strings.ReplaceAll(e.Reason, "\"", "\\\"")
	if e.CausedBy != nil {
		cause := strings.ReplaceAll(e.CausedBy.Error(), "\"", "\\\"")
		return []byte(fmt.Sprintf(`{"type":"%s","reason":"%s","cause":"%s"}`, e.Type, reason, cause)), nil
	}
	return []byte(fmt.Sprintf(`{"type":"%s","reason":"%s"}`, e.Type, reason)), nil
}

func (e *Error) Error() string {
	if e.CausedBy != nil {
		return fmt.Sprintf("type: %s, reason: %s, cause: %s", e.Type, e.Reason, e.CausedBy)
	}
	return fmt.Sprintf("type: %s, reason: %s", e.Type, e.Reason)
}
