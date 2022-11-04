package pgdb

// Hook is executed before a request is maded to the database, the returned
// function is executed after the call is made to the database.
// This can be used to "hook" into the database calls, to perform actions such
// as query logging or request time recording.
type Hook interface {
	Before(query string, args ...interface{}) AfterFunc
}

// AfterFunc is returned by a hook to be called after the request is made to
// the database. Type aliased so that implementors of Hook do not require a
// dependency on this package.
type AfterFunc = func()

// HookFunc is a function that implements the Hook interface so consumers can
// declare an inline hook function rather than creating a new type that
// implements the Hook interface if desired.
type HookFunc func(query string, args ...interface{}) AfterFunc

func (hf HookFunc) Before(query string, args ...interface{}) AfterFunc {
	return hf(query, args...)
}
