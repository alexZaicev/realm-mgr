package pgdb

import (
	"time"
)

const (
	// Istio has a 60m idle timeout at which point it will close any idle
	// connections. This means any long-lived connections, such as those that
	// inhabit a database connection pool, may suddenly stop working, and the
	// only way to find out is to use the closed connection which failed with
	// error messages with terms such as "broken pipe" or "EOF".
	//
	// To fix this, use a timeout smaller than what Istio sets to ensure we
	// clean up idle connections before Istio can close them underneath us. In
	// reality, we can be any value less than 60m, but we use 30m as it's a
	// reasonable balance.
	defaultMaxIdleTimeMinutes = 30
)

// DBOptions can be used to configure the ConnProvider's underlying *sql.DB
// connection pool.
//
// The first 4 options that can be configured using this struct are described here:
// https://golang.org/doc/database/manage-connections
type DBOptions struct {
	MaxIdleTime  *time.Duration
	MaxLifetime  *time.Duration
	MaxOpenConns *int
	MaxIdleConns *int

	Hooks []Hook
}

// DBOption is an interface for applying functional options to a `DBOptions`
// struct.
type DBOption interface {
	apply(*DBOptions)
}

// DBOptionFunc provides a type alias for a function that takes a pointer to a
// `DBOptions` type as an argument.
//
// This allows the type to implement the `DBOption` interface.
type DBOptionFunc func(*DBOptions)

// apply function `f` to `DBOptions`
func (f DBOptionFunc) apply(o *DBOptions) {
	f(o)
}

// WithMaxIdleTime is used to configure the maximum time that a database
// connection can remain idle:
// https://pkg.go.dev/database/sql#DB.SetConnMaxIdleTime
//
// Defaults to 30 minutes to ensure idle connections are cleaned up internally
// before Istio's 60 minute timeout closes them without the database connection
// pool being aware.
func WithMaxIdleTime(maxIdleTime time.Duration) DBOption {
	return DBOptionFunc(func(o *DBOptions) {
		o.MaxIdleTime = &maxIdleTime
	})
}

// WithMaxLifetime is used to configure the maximum time that a database
// connection can be re-used:
// https://pkg.go.dev/database/sql#DB.SetConnMaxLifetime
func WithMaxLifetime(maxLifetime time.Duration) DBOption {
	return DBOptionFunc(func(o *DBOptions) {
		o.MaxLifetime = &maxLifetime
	})
}

// WithMaxOpenConns is used to configure the maximum number of open database
// connections to be opened at any one time:
// https://pkg.go.dev/database/sql#DB.SetMaxOpenConns
func WithMaxOpenConns(maxOpenConns int) DBOption {
	return DBOptionFunc(func(o *DBOptions) {
		o.MaxOpenConns = &maxOpenConns
	})
}

// WithMaxIdleConns is used to configure the maximum number of idling
// connections to be kept open at any one time:
// https://pkg.go.dev/database/sql#DB.SetMaxIdleConns
func WithMaxIdleConns(maxIdleConns int) DBOption {
	return DBOptionFunc(func(o *DBOptions) {
		o.MaxIdleConns = &maxIdleConns
	})
}

// WithHooks is used to add hooks into the database requests
// recorder, e.g. to supply metrics or log queries
// Each hook's Before method is called in the order the hooks
// are provided here, and the After methods are called in FILO order.
// E.g. it would be recommend to register a db request timing hook
// last as then the only operation between its Before and After methods
// will be the actual database call.
func WithHooks(hooks []Hook) DBOption {
	return DBOptionFunc(func(o *DBOptions) {
		o.Hooks = hooks
	})
}
