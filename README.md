# pipe
A simple Go stream processing library that works like Unix pipe.
This library has no external dependencies and is full asynchronous.

To install :
```sh
go get github.com/hyperboloide/pipe
```

### Example

Bellow is a very basic example that:

1. Open a file
2. Compress it
3. Save it

```go
package main

import (
    "compress/gzip"
    "github.com/hyperboloide/pipe"
    "io"
    "log"
    "os"
)

func zip(r io.Reader, w io.Writer) error {
    gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
    if err != nil {
        return err
    }
    defer gzw.Close()
    _, err = io.Copy(gzw, r)
    return err
}

func main() {
    // pipe input
    in, err := os.Open("test.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer in.Close()

    // create a new pipe with a io.Reader
    p := pipe.New(in)

    // Pushing transformation function
    p.Push(zip)

    // pipe output
    out, err := os.Create("test.txt.tgz")
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()

    // Set pipe output io.Writer
    p.To(out)

    // Wait for pipe process to complete
    if err := p.Exec(); err != nil {
        log.Fatal(err)
    }
}
```