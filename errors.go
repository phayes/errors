package errors

import (
	"errors"
	"fmt"
)

// Default implementation of Error interface
type Error struct {
	err   error
	inner error
	anno  map[string]interface{}
}

// This returns a string with error information, excluding inner errors
func (e *Error) Message() string {
	return e.err.Error()
}

// Get the inner error that is wrapped by this error
func (e *Error) Inner() *Error {
	if e.inner == nil {
		return nil
	} else if innerError, ok := e.inner.(*Error); ok {
		return innerError
	} else {
		return &Error{
			err: e.inner,
		}
	}
}

// Get the base error that forms the basis of the Error - returns a copy of itself without inners
func (e *Error) Base() *Error {
	if e.inner == nil {
		return e
	} else {
		return &Error{
			err:  e.err,
			anno: e.anno,
		}
	}
}

// This returns a string with all available error information, including inner
// errors that are wrapped by this errors and annotations.
func (e *Error) Error() string {
	var msg string
	if e.inner != nil {
		msg = e.Message() + ". " + e.inner.Error()
	} else {
		msg = e.Message()
	}
	anno := e.GetAnnotations()
	if len(anno) != 0 {
		msg += ". "
		for key, value := range anno {
			msg += key + ": " + fmt.Sprintf("%s", value) + ". "
		}
	}
	return msg
}

// Annotate additional information about the error
func (e *Error) Annotate(key string, value interface{}) {
	if e.anno == nil {
		e.anno = make(map[string]interface{})
	}
	e.anno[key] = value
}

// Get all annotations about the error. This includes annotations from inner errors.
// If inner errors and outer errors share the same key, the outer error's information takes precedence.
func (e *Error) GetAnnotations() map[string]interface{} {
	var anno map[string]interface{}
	if e.anno == nil {
		anno = make(map[string]interface{})
	} else {
		anno = e.anno
	}

	if innerError := e.Inner(); innerError != nil {
		innerAnno := innerError.GetAnnotations()
		for key, value := range innerAnno {
			if _, ok := anno[key]; !ok {
				anno[key] = value
			}
		}
	}
	return anno
}

// Get the innermost error that is wrapped by this error
func (e *Error) Cause() error {
	if e.Inner() == nil {
		return e
	} else {
		return e.Inner().Cause()
	}
}

// Wrap the passed error in this error and return the new error
func (e *Error) Wrap(err error) *Error {
	e.inner = err
	return e
}

// Wrap the passed error in this error
func (e *Error) Wraps(str string) *Error {
	return e.Wrap(errors.New(str))
}

// Wrap the passed error in this error
func (e *Error) Wrapf(format string, args ...interface{}) *Error {
	return e.Wrap(fmt.Errorf(format, args...))
}

// Append the passed error in this error
func (e *Error) Append(err error) *Error {
	if outerError, ok := err.(*Error); ok {
		return outerError.Wrap(e)
	} else {
		return &Error{
			err:   err,
			inner: e,
		}
	}
}

// Wrap the passed error in this error and return a copy
func (e *Error) Appends(str string) *Error {
	return &Error{
		err:   New(str),
		inner: e,
	}
}

// Wrap the passed error in this error and return a copy
func (e *Error) Appendf(format string, args ...interface{}) *Error {
	return &Error{
		err:   Newf(format, args...),
		inner: e,
	}
}

// Create new error from string.
// It intentionally mirrors the standard "errors" module so as to be a drop-in replacement
func New(s string) *Error {
	return &Error{
		err: errors.New(s),
	}
}

// Same as New, but with fmt.Printf-style parameters.
// This is a replacement for fmt.Errorf.
func Newf(format string, args ...interface{}) *Error {
	return &Error{
		err: errors.New(fmt.Sprintf(format, args...)),
	}
}

// Check to see if two errors are the same
func Equal(e1 error, e2 error) bool {
	if e1 == e2 {
		return true
	}

	// Try to convert them into errors.Error and see if they are the same or if one is based on another
	e1Error, e1ok := e1.(*Error)
	e2Error, e2ok := e2.(*Error)

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

// Check if two errors are the same or if the first contains the second
// This will recursively check their inner components to see if one is an instance of the other
func IsA(outerErr error, innerErr error) bool {
	if Equal(outerErr, innerErr) {
		return true
	}

	// Recursively check to see if the inner is contained in the outer
	if outerError, ok := outerErr.(*Error); ok {
		if outerInner := outerError.Inner(); outerInner != nil {
			return IsA(outerInner, innerErr)
		}
	}

	// No match
	return false
}
