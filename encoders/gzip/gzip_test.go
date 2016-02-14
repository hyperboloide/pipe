package gzip_test

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/encoders/gzip"
	"io/ioutil"
	"os"
	"testing"
)

func TestGzip(t *testing.T) {

	gz := &gzip.Gzip{}
	if err := gz.Start(); err != nil {
		t.Error(err)
	}

	reader, err := os.Open("../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}
	defer reader.Close()

	bin, err := ioutil.ReadFile("../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}

	var result1 bytes.Buffer
	writer := bufio.NewWriter(&result1)

	if err := pipe.New(reader).Push(gz.Encode).To(writer).Exec(); err != nil {
		t.Error(err)
	} else if bytes.Equal(result1.Bytes(), bin) == true {
		t.Error(errors.New("file content match original."))
	} else if len(result1.Bytes()) >= len(bin) {
		t.Error(errors.New("file content greater than original."))
	}

	reader2 := bytes.NewReader(result1.Bytes())
	var result2 bytes.Buffer
	writer = bufio.NewWriter(&result2)

	if err := pipe.New(reader2).Push(gz.Decode).To(writer).Exec(); err != nil {
		t.Error(err)
	} else if bytes.Equal(result1.Bytes(), bin) == true {
		t.Error(errors.New("file content match original."))
	} else if len(result1.Bytes()) >= len(bin) {
		t.Error(errors.New("file content greater than original."))
	}

}
