package errors

import (
	"errors"
	"strings"
)

type baseError struct {
	Err error
	msg string
}

func newBaseError(msg string, err error) baseError {
	return baseError{
		msg: msg,
		Err: err,
	}
}

// Error allows baseError and any structs that embed it to satisfy the error
// interface.
func (e *baseError) Error() string {
	return e.msg
}

// Unwrap allows baseError and any structs that embed it to be used with the
// error wrapping utilities introduced in go 1.13.
func (e *baseError) Unwrap() error {
	// This nil check accounts for the situation where the embedded *baseError
	// in one of the public errors is nil - if it has been constructed without
	// using one of the helper functions (e.g in other package's unit tests).
	if e == nil {
		return nil
	}
	return e.Err
}

// ErrInfoExtractor defines a function signature for functions
// that can extract information from an error.
type ErrInfoExtractor func(error) string

// UnwrapInfoExtractor creates an ErrInfoExtractor function that unwraps an error
// to the specified depth, combining all the messages together into one string.
func UnwrapInfoExtractor(maxDepth int) ErrInfoExtractor {
	return func(err error) string {
		if err == nil {
			return ""
		}

		builder := strings.Builder{}
		builder.WriteString(err.Error())

		for i := 1; i < maxDepth; i++ {
			err = errors.Unwrap(err)
			if err == nil {
				return builder.String()
			}

			builder.WriteString(": " + err.Error())
		}

		return builder.String()
	}
}
