package std

import (
	"errors"
	"io"
	"os"
)

type Stdout struct {
}

func (s *Stdout) Start() error {
	if os.Stdout == nil {
		return errors.New("stdout is not available")
	}
	return nil
}

func (s *Stdout) NewWriter(id string) (io.WriteCloser, error) {
	return os.Stdout, nil
}
