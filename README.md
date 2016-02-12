# pipe

[![Build Status](https://travis-ci.org/hyperboloide/pipe.svg)](https://travis-ci.org/hyperboloide/pipe)
[![GoDoc](https://godoc.org/github.com/hyperboloide/pipe?status.svg)](https://godoc.org/github.com/hyperboloide/pipe)

A simple stream processing library that works like Unix pipes.
This library has no external dependencies and is fully asynchronous.
Create a Pipe from a Reader, add some transformation functions and get the result writed to a Writer.

To install :
```sh
$ go get github.com/hyperboloide/pipe
```

Then add the following import :
```go
import "github.com/hyperboloide/pipe"
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

    // pipe output
    out, err := os.Create("test.txt.tgz")
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()

    // create a new pipe with a io.Reader
    // Push a transformation function
    // Set output
    // Exec and get errors if any
    if err := pipe.New(in).Push(zip).To(out).Exec(); err != nil {
        log.Fatal(err)
    }
}
```
