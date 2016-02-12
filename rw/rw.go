package rw

import (
	"io"
)

// an interface that defines a start function, used for setup.
type Base interface {
	Start() error
}

// a base interface to write data.
type Writer interface {
	Base
	NewWriter(string) (io.WriteCloser, error)
}

// A base interface to read data.
type Reader interface {
	Base
}

// A base interface to delete data.
type Deleter interface {
	Base
	Delete(string) error
}

// A type with Base, Reader and Writer
type ReadWriter interface {
	Base
	NewReader(string) (io.ReadCloser, error)
	NewWriter(string) (io.WriteCloser, error)
}

// A type with Base, Reader, Writer and Deleter
type ReadWriteDeleter interface {
	Base
	NewReader(string) (io.ReadCloser, error)
	NewWriter(string) (io.WriteCloser, error)
	Delete(string) error
}

// The Prefix struct allows to define prefix and suffix (for example a
// file extension)
type Prefixed struct {
	Prefix string
	Suffix string
}

// Generate a name from with prefix and suffix
func (p *Prefixed) Name(id string) string {
	return p.Prefix + id + p.Suffix
}
