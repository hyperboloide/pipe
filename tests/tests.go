package tests

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/encoders"
	"github.com/hyperboloide/pipe/rw"
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

func TestEncoderDecoder(ed encoders.EncoderDecoder, file string) error {
	var err error
	var encoded bytes.Buffer
	encodedWriter := bufio.NewWriter(&encoded)

	var decoded bytes.Buffer
	decodedWriter := bufio.NewWriter(&decoded)

	var originalReader *os.File
	var originalByte []byte

	// Should start without errors
	if err := ed.Start(); err != nil {
		return err
	}

	// should open the test file and return a reader
	if originalReader, err = os.Open(file); err != nil {
		return err
	} else {
		defer originalReader.Close()
	}

	// should open the test file and return the bytes
	if originalByte, err = ioutil.ReadFile("../../tests/test.jpg"); err != nil {
		return err
	}

	// should encode test file reader
	if err := pipe.New(originalReader).Push(ed.Encode).To(encodedWriter).Exec(); err != nil {
		return err
	}

	// the encoded bytes should not match the original
	if bytes.Equal(encoded.Bytes(), originalByte) == true {
		return errors.New("the encoded bytes should not match the original")

	}

	encodedReader := bytes.NewReader(encoded.Bytes())
	// should decoded the encoded writer
	if err := pipe.New(encodedReader).Push(ed.Decode).To(decodedWriter).Exec(); err != nil {
		return err
	}

	// the decoded bytes should match the original
	if bytes.Equal(decoded.Bytes(), originalByte) == false {
		return errors.New("the decoded bytes should match the original")

	}

	return nil
}
