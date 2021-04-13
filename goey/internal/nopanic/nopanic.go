package nopanic

import (
	"fmt"
	"runtime/debug"
)

// PanicError represents a panic that occurred.
type PanicError struct {
	value interface{}
	stack string
}

// Error returns a description of the error.
func (pe PanicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", pe.value, pe.stack)
}

// Value returns the value returned by recover after a panic.
func (pe PanicError) Value() interface{} {
	return pe.value
}

// Stack returns a formatted stack trace of the goroutine that originally
// paniced.
func (pe PanicError) Stack() string {
	return pe.stack
}

// New wraps the value returned by recover into a PanicError.
func New(value interface{}) PanicError {
	stack := string(debug.Stack())
	return PanicError{value, stack}
}

// Wrap ensures that no panics escape.  Action will be called, and if it
// returns normally, its return value will be returned.  However, if action
// panics, the panic will be converted to an error, and will be returned.
func Wrap(action func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = New(r)

			// If r is an instance of PanicError, we should consider using it
			// as is, rather than adding another layer of wrapping.
		}
	}()

	return action()
}

// Unwrap will panic if err is an instance of PanicError, otherwise it will
// return the error unmodified.
func Unwrap(err error) error {
	if pe, ok := err.(PanicError); ok {
		panic(pe)
	}

	return err
}
