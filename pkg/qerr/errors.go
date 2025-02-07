package qerr

import (
	"fmt"
	"runtime/debug"
)

type ErrorCode int

const (
	Unknown      ErrorCode = iota
	InvalidInput           // For validation errors
	NotFound               // For missing resources
	Internal               // For system errors
)

type QError struct {
	Inner      error
	Message    string
	StackTrace string
	Code       ErrorCode
	Misc       map[string]any
}

// WrapError creates a QError. It adds more information to the error
// for debugging.
func Wrap(err error, code ErrorCode, messagef string, msgArgs ...any) QError {
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
