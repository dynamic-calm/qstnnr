package qstnnr

import (
	"fmt"
	"runtime/debug"
)

type Error struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]any
}

func wrapError(err error, messagef string, msgArgs ...any) Error {
	return Error{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}

func (err Error) Error() string {
	return err.Message
}
