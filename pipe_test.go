//
// pipe_test.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Apr 27 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package pipe_test

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"testing"
	"github.com/hyperboloide/pipe"
	"os"
	"log"
)

func genBlob(size int) []byte {
	blob := make([]byte, size)
	for i := 0; i < size; i++ {
		blob[i] = 65 // ascii 'A'
	}
	return blob
}

var bin = genBlob(1 << 24)

func passProc(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}

func zip(r io.Reader, w io.Writer) error {
	gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return err
	}
	defer gzw.Close()
	_, err = io.Copy(gzw, r)
	return err
}

func unzip(r io.Reader, w io.Writer) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()
	_, err = io.Copy(w, gzr)
	return err
}

func TestBasic(t *testing.T) {
	binReader := bytes.NewReader(bin)

	var result bytes.Buffer
	writer := bufio.NewWriter(&result)

	p := pipe.New(binReader).To(writer)

	if err := p.Exec(); err != nil {
		t.Errorf("errors detected during pipe: %s", err)
	}

	if !bytes.Equal(result.Bytes(), bin) {
		t.Errorf("result do not match")
	}

	if p.TotalIn != int64(len(bin)) {
		t.Errorf("TotalIn do not match")
	}

	if p.TotalOut != p.TotalIn {
		t.Errorf("TotalOut: %d should equal TotalIn", p.TotalOut)
	}
}

func TestProcess(t *testing.T) {
	p := pipe.New(bytes.NewReader(bin)).Push(passProc, zip, unzip, zip)

	var result bytes.Buffer
	writer := bufio.NewWriter(&result)
	p.To(writer)

	if err := p.Exec(); err != nil {
		t.Errorf("errors detected during pipe: %s", err)
	}

	if bytes.Equal(result.Bytes(), bin) {
		t.Errorf("result should not match")
	}
}

func TestError(t *testing.T) {
	p := pipe.New(bytes.NewReader(bin))

	var procErr = func(r io.Reader, w io.Writer) error {
		io.Copy(w, r)
		return errors.New("some error!")
	}

	p.Push(passProc, passProc, procErr, passProc, passProc)

	var result bytes.Buffer
	writer := bufio.NewWriter(&result)
	p.To(writer)

	if err := p.Exec(); err == nil {
		t.Errorf("pipe should have an error")
	}
}

func TestTee(t *testing.T) {
	p := pipe.New(bytes.NewReader(bin))
	p.Push(zip)

	pTee := p.Tee()

	p.Push(unzip)
	pTee.Push(unzip)

	var result bytes.Buffer
	writer := bufio.NewWriter(&result)
	p.To(writer)

	var resultTee bytes.Buffer
	writerTee := bufio.NewWriter(&resultTee)
	pTee.To(writerTee)

	if err := p.Exec(); err != nil {
		t.Errorf("pipe should not have error %s", err)
	}

	if err := pTee.Exec(); err != nil {
		t.Errorf("pipe should not have error %s", err)
	}

	if !bytes.Equal(result.Bytes(), bin) {
		t.Errorf("result do not match")
	}
	if !bytes.Equal(resultTee.Bytes(), bin) {
		t.Errorf("result of tee do not match")
	}
}

func ExamplePipe() {

	// some PipeFilter transformation function
	var zip = func(r io.Reader, w io.Writer) error {
		gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			return err
		}
		defer gzw.Close()
		_, err = io.Copy(gzw, r)
		return err
	}

	// pipe input
	in, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	// create a new pipe with a io.Reader
	p := pipe.New(in)

	// Pushing transformation function
	p.Push(zip)

	// pipe output
	out, err := os.Create("test.txt.tgz")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Set pipe output io.Writer
	p.To(out)

	// Wait for pipe process to complete
	if err := p.Exec(); err != nil {
		log.Fatal(err)
	}

}
