package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

func Test_NewInternalError_Success(t *testing.T) {
	err := realmmgr_errors.NewInternalError("hello world", errors.New("mock error"))
	assert.EqualError(t, err, "an internal error occurred: hello world")
	assert.IsType(t, realmmgr_errors.InternalErrorType, err)
	assert.EqualError(t, errors.Unwrap(err), "mock error")
}

func Test_NewInvalidArgumentError_Success(t *testing.T) {
	err := realmmgr_errors.NewInvalidArgumentError("arg1", "hello world")
	assert.EqualError(t, err, "an invalid argument error occurred: argument arg1 hello world")
	assert.IsType(t, realmmgr_errors.InvalidArgumentErrorType, err)
	assert.NoError(t, errors.Unwrap(err))
}

func Test_NewUnknownError_Success(t *testing.T) {
	err := realmmgr_errors.NewUnknownError("hello world", errors.New("mock error"))
	assert.EqualError(t, err, "an unknown error occurred: hello world")
	assert.IsType(t, realmmgr_errors.UnknownErrorType, err)
	assert.EqualError(t, errors.Unwrap(err), "mock error")
}

func Test_NewNotFoundError_Success(t *testing.T) {
	err := realmmgr_errors.NewNotFoundError("hello world", errors.New("mock error"))
	assert.EqualError(t, err, "not found error occurred: hello world")
	assert.IsType(t, realmmgr_errors.NotFoundErrorType, err)
	assert.EqualError(t, errors.Unwrap(err), "mock error")
}
