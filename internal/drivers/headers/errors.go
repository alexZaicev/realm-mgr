package headers

import (
	"fmt"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

var (
	HeaderNotFoundType       = &HeaderNotFound{}
	MultipleHeadersFoundType = &MultipleHeadersFound{}
)

// HeaderNotFound indicates that a header was expected in a request but was not
// present.
type HeaderNotFound struct {
	*realmmgr_errors.NotFoundError
	Name string
}

// NewHeaderNotFound creates a new HeaderNotFound.
func NewHeaderNotFound(name string) *HeaderNotFound {
	return &HeaderNotFound{
		NotFoundError: realmmgr_errors.NewNotFoundError(fmt.Sprintf("header %s not found", name), nil),
		Name:          name,
	}
}

// PresentableError returns an error message appropriate to output to consumers of a service.
func (e *HeaderNotFound) PresentableError() string {
	return fmt.Sprintf(`Header %q was not found`, e.Name)
}

// MultipleHeadersFound indicates that multiple header values were found in a
// request where only one value was expected.
type MultipleHeadersFound struct {
	*realmmgr_errors.InvalidArgumentError
	Name string
}

// NewMultipleHeadersFound creates a new MultipleHeadersFound.
func NewMultipleHeadersFound(name string) *MultipleHeadersFound {
	return &MultipleHeadersFound{
		InvalidArgumentError: realmmgr_errors.NewInvalidArgumentError(
			fmt.Sprintf("header %q", name),
			"cannot have multiple values",
		),
		Name: name,
	}
}

// PresentableError returns an error message appropriate to output to consumers of a service.
func (e *MultipleHeadersFound) PresentableError() string {
	return fmt.Sprintf(`Header %q cannot have multiple values`, e.Name)
}
