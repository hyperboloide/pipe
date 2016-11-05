package encoders

import (
	"io"
)

// BaseEncoder is an interface that defines a start function, used for setup.
type BaseEncoder interface {
	Start() error
}

// EncodeFun is a shortcut type for an encode function
type EncodeFun func(io.Reader, io.Writer) error
