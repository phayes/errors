package errors

import (
	"fmt"
)

type Error interface {
	// This returns the error message without inner errors
	Message() string

	// Get the inner error. Will return nil if there is no inner error
	Inner() error

	// Wrap the given error. Calling Inner() should retreive it. Return a copy of the outer error as an Error.
	Wrap(error) Error

	// If we are using an error from another package as the outer (wrapping) error, we can get it here
	// If this Error is not derived from another error, Base() should return a copy itself without any inners
	Base() error

	// Implements the built-in error interface.
	Error() string
}

// Default implementation of Error interface
type DefaultError struct {
	msg   string
	inner error
	base  error
}

// This returns a string with error information, excluding inner errors
func (e DefaultError) Message() string {
	return e.msg
}

// This returns a string with all available error information, including inner
// errors that are wrapped by this errors.
func (e DefaultError) Error() string {
	if e.inner != nil {
		return e.Message() + " " + e.inner.Error()
	} else {
		return e.Message()
	}
}

// Get the inner error that is wrapped by this error
func (e DefaultError) Inner() error {
	return e.inner
}

// Get the base error that forms the basis of the DefaultError
// If this error was not constructed out of a different one then Base() returns a copy of itself without inners
func (e DefaultError) Base() error {
	if e.base != nil {
		return e.base
	} else {
		e.inner = nil
		return e
	}
}

// Wrap the passed error in this error and return a copy
func (e DefaultError) Wrap(err error) Error {
	e.inner = err
	return e
}

// Create new error from string.
// It intentionally mirrors the standard "errors" module so as to be a drop-in replacement
func New(s string) error {
	return DefaultError{
		msg: s,
	}
}

// Same as New, but with fmt.Printf-style parameters.
// This is a replacement for fmt.Errorf.
func Newf(format string, args ...interface{}) error {
	return DefaultError{
		msg: fmt.Sprintf(format, args...),
	}
}

// Wrap the first error in the second error.
func Wrap(innerErr error, outerErr error) error {
	if outerError, ok := outerErr.(Error); ok {
		outerError = outerError.Wrap(innerErr)
		return outerError
	} else {
		return DefaultError{
			msg:   outerErr.Error(),
			inner: innerErr,
			base:  outerErr,
		}
	}
}

// Wrap an error in a new error using the provided string
// Functionally similar to using errors.New then errors.Wrap
func Wraps(err error, s string) error {
	return DefaultError{
		msg:   s,
		inner: err,
	}
}

// Same as Wraps, but with fmt.Printf-style parameters.
func Wrapf(err error, format string, args ...interface{}) error {
	return DefaultError{
		msg:   fmt.Sprintf(format, args...),
		inner: err,
	}
}

// Check to see if two errors are the same
func Equal(e1 error, e2 error) bool {
	if e1 == e2 {
		return true
	}

	// Try to convert them into errors.Error and see if they are the same or if one is based on another
	e1Error, e1ok := e1.(Error)
	e2Error, e2ok := e2.(Error)

	// If neither of them are Errors we can stop now
	if !e1ok && !e2ok {
		return false
	}
	if e1ok && e2ok {
		if e1Error == e2Error || e2Error.Base() == e1Error.Base() {
			return true
		}
	}
	if e1ok {
		if e1Error.Base() == e2 {
			return true
		}
	}
	if e2ok {
		if e2Error.Base() == e1 {
			return true
		}
	}

	// No match
	return false
}

// Check if two errors are the same or if the first contains the second
// This will recursively check their inner components to see if one is an instance of the other
func IsA(outerErr error, innerErr error) bool {
	if Equal(outerErr, innerErr) {
		return true
	}

	// Recursively check to see if the inner is somehow contained in the outer
	if outerError, ok := outerErr.(Error); ok {
		if outerInner := outerError.Inner(); outerInner != nil {
			return IsA(outerInner, innerErr)
		}
	}

	// No match
	return false
}
