package qstnnr

import (
	"fmt"
	"runtime/debug"
)

type ErrorCode int

const (
	ErrorCodeUnknown      ErrorCode = iota
	ErrorCodeInvalidInput           // For validation errors
	ErrorCodeNotFound               // For missing resources
	ErrorCodeInternal               // For system errors
)

type QError struct {
	Inner      error
	Message    string
	StackTrace string
	Code       ErrorCode
	Misc       map[string]any
}

// wrapError creates a QError. It adds more information to the error
// for debugging.
func wrapError(err error, code ErrorCode, messagef string, msgArgs ...any) QError {
	return QError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Code:       code,
		Misc:       make(map[string]any),
	}
}

func (err QError) Error() string {
	return err.Message
}
