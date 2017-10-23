package gzip_test

import (
	"testing"

	"github.com/hyperboloide/pipe/encoders/gzip"
	"github.com/hyperboloide/pipe/tests"
)

func TestGzip(t *testing.T) {

	enc := &gzip.Gzip{}
	if err := enc.Start(); err != nil {
		t.Error(err)
	}

	if err := tests.TestEncoderDecoder(enc, "../../tests/test.jpg"); err != nil {
		t.Error(err)
	}

}
