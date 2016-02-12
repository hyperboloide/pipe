package rw_test

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/rw"
)

func genBlob(size int) []byte {
	blob := make([]byte, size)
	for i := 0; i < size; i++ {
		blob[i] = 65 // ascii 'A'
	}
	return blob
}

var bin = genBlob(1 << 24)

func testReadWriteDeleter(rwd rw.ReadWriteDeleter, id string) error {

	if err := rwd.Start(); err != nil {
		return err
	}

	w, err := rwd.NewWriter(id)
	if err != nil {
		return err
	}

	binReader := bytes.NewReader(bin)

	if err := pipe.New(binReader).ToCloser(w).Exec(); err != nil {
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
	} else if bytes.Equal(result.Bytes(), bin) == false {
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
