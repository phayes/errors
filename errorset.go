package errors

import (
	"strings"
	"sync"
)

// ErrorSet is a set of errors that can be collected together in a non-heirarchical manner.
type ErrorSet struct {
	mu  sync.Mutex
	set map[string]error
}

// NewErrorSet creates a new empty ErrorSet
func NewErrorSet() *ErrorSet {
	return &ErrorSet{}
}

// Return Error as string
func (e *ErrorSet) Error() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e == nil || e.set == nil || len(e.set) == 0 {
		return ""
	}
	output := ""
	for str, err := range e.set {
		output += str + ": " + err.Error() + ". "
	}
	output = strings.TrimRight(output, ". ")
	return output
}

// Add an error to the error set
func (e *ErrorSet) Add(key string, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.set == nil {
		e.set = map[string]error{}
	}

	e.set[key] = err
}

// Get an error from the error set
func (e *ErrorSet) Get(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.set == nil {
		return nil
	}
	return e.set[key]
}

// GetAll Gets all errors from the error set
func (e *ErrorSet) GetAll() map[string]error {
	e.mu.Lock()
	defer e.mu.Unlock()

	set := map[string]error{}

	if e.set == nil {
		return set
	}

	for k, v := range e.set {
		set[k] = v
	}

	return set
}

// HasErrors returns true if the ErrorSet contains an error
func (e *ErrorSet) HasErrors() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.set == nil || len(e.set) == 0 {
		return false
	}

	return true
}
