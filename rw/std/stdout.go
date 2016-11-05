package std

import (
	"errors"
	"io"
	"os"
)

// Stdout allows for writing on the standard output.
type Stdout struct {
}

// Start by checking that stdout is available.
func (s *Stdout) Start() error {
	if os.Stdout == nil {
		return errors.New("stdout is not available")
	}
	return nil
}

// NewWriter returns a new writer tot stdout
func (s *Stdout) NewWriter(id string) (io.WriteCloser, error) {
	return os.Stdout, nil
}
