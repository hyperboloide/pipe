package tests

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/rw"
	"io/ioutil"
	"os"
)

func TestReadWriteDeleter(rwd rw.ReadWriteDeleter, id, file string) error {

	if err := rwd.Start(); err != nil {
		return err
	}

	w, err := rwd.NewWriter(id)
	if err != nil {
		return err
	}

	reader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := pipe.New(reader).ToCloser(w).Exec(); err != nil {
		return err
	}

	r, err := rwd.NewReader(id)
	if err != nil {
		return err
	}

	var result bytes.Buffer
	writer := bufio.NewWriter(&result)

	if err := pipe.New(r).To(writer).Exec(); err != nil {
		return err
	}
	bin, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if bytes.Equal(result.Bytes(), bin) == false {
		return errors.New("file content do not match original.")
	}

	if err := rwd.Delete(id); err != nil {
		return err
	}

	if _, err := rwd.NewReader(id); err == nil {
		return errors.New("object not deleted.")
	}
	return nil
}
