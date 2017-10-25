package rw

import (
	"io"
)

// Base is an interface that defines a start function, used for setup.
type Base interface {
	Start() error
}

// Writer is a base interface to write data.
type Writer interface {
	Base
	NewWriter(string) (io.WriteCloser, error)
}

// Reader is a base interface to read data.
type Reader interface {
	Base
	NewReader(string) (io.ReadCloser, error)
}

// Deleter is a base interface to delete data.
type Deleter interface {
	Base
	Delete(string) error
}

// ReadWriter is a type with Base, Reader and Writer
type ReadWriter interface {
	Base
	NewReader(string) (io.ReadCloser, error)
	NewWriter(string) (io.WriteCloser, error)
}

// ReadWriteDeleter is a type with Base, Reader, Writer and Deleter
type ReadWriteDeleter interface {
	Base
	NewReader(string) (io.ReadCloser, error)
	NewWriter(string) (io.WriteCloser, error)
	Delete(string) error
}

// ReadDeleter is a type with Base, Reader and Deleter
type ReaderDeleter interface {
	Base
	NewReader(string) (io.ReadCloser, error)
	Delete(string) error
}

// WriteDeleter is a type with Base, Writer and Deleter
type WriterDeleter interface {
	Base
	NewWriter(string) (io.WriteCloser, error)
	Delete(string) error
}

// Prefixed  struct allows to define prefix and suffix (for example a
// file extension)
type Prefixed struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

// Name generate a name from with prefix and suffix
func (p *Prefixed) Name(id string) string {
	return p.Prefix + id + p.Suffix
}
