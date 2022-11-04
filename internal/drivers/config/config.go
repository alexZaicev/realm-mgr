package config

import (
	"fmt"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

// Config defines a minimal interface required for a config implementation
// to provide. Rather than handle the nil checking and type casting explicitly,
// it is recommended to use the generic GetConfig and MustConfig functions.
type Config interface {
	Get(string) any
	Keys() []string
}

// Get provides typed access to config values in the provided Config
// implementation. If a config entry with the specified key does not exist
// a config.NotFound error is returned. If the entry exists but the value is
// of a different type to that requested a config.ValueNotExpectedType error
// is returned.
func Get[T any](c Config, key string) (T, error) {
	var zeroVal T
	val := c.Get(key)
	if val == nil {
		return zeroVal, realmmgr_errors.NewNotFoundError(
			fmt.Sprintf("configuration for key %s not found", key),
			nil,
		)
	}

	castVal, ok := val.(T)
	if !ok {
		return zeroVal, realmmgr_errors.NewInvalidArgumentError(
			key,
			fmt.Sprintf("value is of type %T, not %T", zeroVal, val),
		)
	}

	return castVal, nil
}

// Must mirrors the logic of Get but panics on error rather
// than returning the error.
func Must[T any](c Config, key string) T {
	val, err := Get[T](c, key)
	if err != nil {
		panic(err)
	}

	return val
}
