/*
Go's standard `errors` package is a simple error package for handling errors, and works fine for shallow codebases, but has several shortcomings that `phayes/errors`
hopes to remediate.

The biggest anti-pattern that the standard errors package encourages is one we are all familiar with:

    if err != nil
        return errors.New(err.Error() + ". Some more details about the error")
    }


This anti-pattern is corrected in phayes/errors by allowing you to cleanly wrap one error with another like so:

    if err != nil {
        return errors.Append(err, "Some more details about the error")
    }


This allows us to cleanly add more details to an error, while preseving the underlying error for later inspection.

It also plays nicely with standard library errors and predefined errors:

    import (
        "github.com/phayes/errors"
        "io"
    )

    var ErrFailedStream = errors.New("Failed to read stream.")

    func ReadStream(b []byte) error {
        n, err := reader.Read(b)
        if err != nil {
            return errors.Append(ErrFailedStream, err)
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


Wrapping errors

At it's most basic, `phayes/errors` is a drop in replacement for the standard error package.

    err := errors.New("Could not parse input")

However, it also provides the ability to wrap an error to give it more context:

    import (
        "github.com/phayes/errors"
    )

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


Inspecting errors=

Use the `IsA` function to check to if the error, or any of it's inner errors, is what you're after. This is fully compatible with errors that
are not part of phayes/errors. For example:

    func main() {
        var b []byte
        err := ReadStream(b)
        if err != nil {
            if errors.IsA(err, io.EOF) {
                return // Success!
            } else {
                log.Fatal(err)
            }
        }
    }
*/
package errors
