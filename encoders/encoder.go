package encoders

import (
	"io"
)

// an interface that defines a start function, used for setup.
type BaseEncoder interface {
	Start() error
}

type EncodeFun func(io.Reader, io.Writer) error
