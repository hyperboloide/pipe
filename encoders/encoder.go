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

// Encoder is a BaseEncoder with and Encode method
type Encoder interface {
	Encode(r io.Reader, w io.Writer) error
	BaseEncoder
}

// Decoder is a BaseEncoder with and Decode method
type Decoder interface {
	Decode(r io.Reader, w io.Writer) error
	BaseEncoder
}

// EncoderDecoder is a BaseEncoder with a Encode and Decode methods
type EncoderDecoder interface {
	Encode(r io.Reader, w io.Writer) error
	Decode(r io.Reader, w io.Writer) error
	BaseEncoder
}
