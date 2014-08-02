At it's most basic, it is a drop in replacement for the standard error package.

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
	        return errors.Wrap(err,ErrFailedStream)
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