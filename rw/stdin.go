package rw

import (
	"errors"
	"io"
	"os"
)

type Stdin struct {
}

func (s *Stdin) Start() error {
	if os.Stdin == nil {
		return errors.New("stdin is not available")
	}
	return nil
}

func (s *Stdin) NewReader(id string) (io.ReadCloser, error) {
	return os.Stdin, nil
}
