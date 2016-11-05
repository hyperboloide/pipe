package std

import (
	"errors"
	"io"
	"os"
)

// Stdin allows for reading from the standard input.
type Stdin struct {
}

// Start the readear by checking that stdin is available
func (s *Stdin) Start() error {
	if os.Stdin == nil {
		return errors.New("stdin is not available")
	}
	return nil
}

// NewReader returns a new reader from stdin
func (s *Stdin) NewReader(id string) (io.ReadCloser, error) {
	return os.Stdin, nil
}
