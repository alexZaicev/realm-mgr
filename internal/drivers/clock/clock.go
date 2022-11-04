package clock

import "time"

// Clock defines an interface for getting the current time
// and getting the duration since/until a time.
type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
	Until(t time.Time) time.Duration
}

// StdLibClock provides a real implementation of the Clock interface
// that uses the std library time module.
type StdLibClock struct{}

// NewStdLibClock returns a StdLibClock, and can be used in dependency
// injection frameworks.
func NewStdLibClock() StdLibClock {
	return StdLibClock{}
}

// Now returns the current time in the UTC timezone using time.Now(). If an alternative timezone is
// required, it should be applied by the caller.
func (s StdLibClock) Now() time.Time {
	return time.Now().UTC()
}

// Since returns the duration since the given time using time.Since().
func (s StdLibClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Until returns the duration until the given time using time.Until().
func (s StdLibClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}
