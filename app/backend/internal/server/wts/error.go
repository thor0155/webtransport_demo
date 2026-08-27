package wts

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

type errorTerminates struct {
	error
}

// var ErrTerminates

func WithErrorTerminates(err error) error {
	return &errorTerminates{err}
}

func (e *errorTerminates) Unwrap() error { return e.error }

func (e *errorTerminates) Format(s fmt.State, v rune) {
	errors.FormatError(e, s, v)
}
