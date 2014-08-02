
phayes/errors: Better error handling for go
-----------------------------------------------
[![Build Status](https://travis-ci.org/phayes/errors.svg?branch=master)](https://travis-ci.org/phayes/errors)
[![Coverage Status](https://coveralls.io/repos/phayes/errors/badge.png?branch=master)](https://coveralls.io/r/phayes/errors)

Documentation: https://godoc.org/github.com/phayes/errors

Go standard `errors` package is a great simple error package for producing errors, but has several shortcomings that `phayes/errors`
hopes to remediate. 

The most glaring anti-pattern that the standard `errors` package unintentionally encourages is 
`return errors.New(err.Error() + ' Some more details about the error')`. This anti-pattern is corrected in `phayes/errors` by allowing
you to cleanly wrap one error with another like so: `return errors.Wrap(err, 'More details about the error')` while preseving the
underlying error for later inspection. 

This package also provides several other useful functions for encouraging good patterns in general error handling.

Examples
--------

At it's most basic, `phayes/errors` is a drop in replacement for the standard error package.

```go
    err := errors.New("Could not parse input")
```

However, it also provides the ability to wrap an error to give it more context

```go
    func ReadStream(b []byte) error {
        n, err := reader.Read(b)
        if err != nil {
	        return errors.Wraps(err, "Failed to read stream.")
        }	
    }

    func main() {
    	var b []byte
    	err := ReadStream(b)
    	if err != nil {
    		log.Fatal(err) // Prints "Failed to read stream. unexpected EOF"
    	}
    }
```

It also plays nicely with standard library errors and predefined errors

```go
    import (
        "io"
    )

    var ErrFailedStream = errors.New("Failed to read stream.")

    func ReadStream(b []byte) error {
        n, err := reader.Read(b)
        if err != nil {
	        return errors.Wrap(err, ErrFailedStream)
        }	
    }

    func main() {
    	var b []byte
    	err := ReadStream(b)
    	if err != nil {
    		if errors.IsA(err, io.EOF) {
    			return // Success!
    		} else {
    		    log.Fatal(err) // Prints "Failed to read stream. unexpected EOF"
    		}
    	}
    }
```
