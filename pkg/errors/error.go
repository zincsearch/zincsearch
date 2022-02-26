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
