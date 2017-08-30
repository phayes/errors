package errors

import (
	"errors"
	"fmt"
)

// Error interface may be implemented by other packages
type Error interface {
	// This returns the error message without inner errors
	Message() string

	// Get the inner error. Will return nil if there is no inner error
	Inner() error

	// Wrap the given error. Calling Inner() should retreive it. Return a copy of the outer error as an Error.
	Wrap(error) Error

	// Base() should return a copy of the Error without any inners
	// This method is called to check two errors for equality
	Base() error

	// Implements the built-in error interface.
	Error() string
}

// DefaultError is the default implementation of Error interface
type DefaultError struct {
	err   error
	inner error
}

// Message returns a string with error information, excluding inner errors
func (e DefaultError) Message() string {
	return e.err.Error()
}

// Error returns a string with all available error information, including inner
// errors that are wrapped by this errors.
func (e DefaultError) Error() string {
	if e.inner != nil {
		return e.Message() + ". " + e.inner.Error()
	} else {
		return e.Message()
	}
}

// Inner gets the inner error that is wrapped by this error
func (e DefaultError) Inner() error {
	return e.inner
}

// Base gets the base error that forms the basis of the DefaultError - returns a copy of itself without inners
func (e DefaultError) Base() error {
	return e.err
}

// Wrap the passed error in this error and return a copy
func (e DefaultError) Wrap(err error) Error {
	e.inner = err
	return e
}

// New create new error from string.
// It intentionally mirrors the standard "errors" module so as to be a drop-in replacement
func New(s string) error {
	return DefaultError{
		err: errors.New(s),
	}
}

// Newf is the same as New, but with fmt.Printf-style parameters.
// This is a replacement for fmt.Errorf.
func Newf(format string, args ...interface{}) error {
	return DefaultError{
		err: errors.New(fmt.Sprintf(format, args...)),
	}
}

// Append more information to the error. The reverse of Wrap.
func Append(outerErr error, innerErr error) error {
	if outerError, ok := outerErr.(Error); ok {
		return outerError.Wrap(innerErr)
	}
	return DefaultError{
		err:   outerErr,
		inner: innerErr,
	}
}

// Appends more information to the error using a string. The reverse of Wraps.
func Appends(outerErr error, inner string) error {
	if outerError, ok := outerErr.(Error); ok {
		return outerError.Wrap(New(inner))
	}
	return DefaultError{
		err:   outerErr,
		inner: errors.New(inner),
	}
}

// Appendf appends more information to the error using formatting. The reverse of Wrapf.
func Appendf(outerErr error, format string, args ...interface{}) error {
	if outerError, ok := outerErr.(Error); ok {
		return outerError.Wrap(Newf(format, args...))
	}
	return DefaultError{
		err:   outerErr,
		inner: Newf(format, args...),
	}
}

// Wrap the first error in the second error. Reverse of Append
func Wrap(innerErr error, outerErr error) error {
	if outerError, ok := outerErr.(Error); ok {
		return outerError.Wrap(innerErr)
	}
	return DefaultError{
		err:   outerErr,
		inner: innerErr,
	}
}

// Wraps wraps an error in a new error using the provided string. Reverse of Appends
func Wraps(err error, outer string) error {
	return DefaultError{
		err:   errors.New(outer),
		inner: err,
	}
}

// Wrapf is the same as Wraps, but with fmt.Printf-style parameters. Reverse of Appendf
func Wrapf(err error, format string, args ...interface{}) error {
	return DefaultError{
		err:   errors.New(fmt.Sprintf(format, args...)),
		inner: err,
	}
}

// Equal checks to see if two errors are the same
func Equal(e1 error, e2 error) bool {
	if e1 == e2 {
		return true
	}

	// Try to convert them into errors.Error and see if they are the same or if one is based on another
	e1Error, e1ok := e1.(Error)
	e2Error, e2ok := e2.(Error)

	switch {
	case e1ok && e2ok:
		return e2Error.Base() == e1Error.Base()
	case e1ok && !e2ok:
		return e1Error.Base() == e2
	case !e1ok && e2ok:
		return e2Error.Base() == e1
	case !e1ok && !e2ok:
		return e1 == e2
	default:
		return false
	}
}

// IsA checks if two errors are the same or if the first contains the second
// This will recursively check their inner components to see if one is an instance of the other
func IsA(outerErr error, innerErr error) bool {
	if Equal(outerErr, innerErr) {
		return true
	}

	// Recursively check to see if the inner is contained in the outer
	if outerError, ok := outerErr.(Error); ok {
		if outerInner := outerError.Inner(); outerInner != nil {
			return IsA(outerInner, innerErr)
		}
	}

	// No match
	return false
}

// Cause returns the root cause of the given error. If err does not implement phayes.Error, it returns err itself.
func Cause(err error) error {
	outerError, ok := err.(Error)
	if !ok {
		return err
	}

	if outerError.Inner() == nil {
		return outerError.Base()
	}
	return Cause(outerError.Inner())
}
