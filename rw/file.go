package rw

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Prefixed
	Dir         string
	AllowSub    bool
	RemoveEmpty bool
}

func (s *File) Start() error {

	if s.Dir == "" {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		s.Dir = dir
	}
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	return nil
}

func (s *File) NewWriter(id string) (io.WriteCloser, error) {
	name := s.Prefixed.Name(id)

	if filepath.Dir(name) == "." {
		return os.OpenFile(s.join(name), os.O_RDWR|os.O_CREATE, 0600)
	} else if s.AllowSub == false {
		return nil, errors.New("sub directories not allowed")
	}
	if err := os.MkdirAll(s.join(filepath.Dir(name)), 0700); err != nil {
		return nil, err
	}
	return os.OpenFile(s.join(name), os.O_RDWR|os.O_CREATE, 0600)
}

func (s *File) NewReader(id string) (io.ReadCloser, error) {
	return os.OpenFile(s.join(s.Prefixed.Name(id)), os.O_RDONLY, 0400)
}

func (s *File) Delete(id string) error {
	name := s.Prefixed.Name(id)
	if err := os.Remove(s.join(name)); err != nil {
		return err
	}
	if s.RemoveEmpty == true && filepath.Dir(name) != "." {
		return s.removeIfEmpty(filepath.Dir(name))
	}
	return nil
}

func (s *File) join(path string) string {
	return filepath.Join(s.Dir, path)
}

func (s *File) removeIfEmpty(dir string) error {
	if dir == "." {
		return nil
	} else {
		f, err := os.Open(s.join(dir))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Readdir(1)
		if err != io.EOF {
			return err
		}
	}
	os.Remove(s.join(dir))
	return s.removeIfEmpty(filepath.Dir(dir))
}
