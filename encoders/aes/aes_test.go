package aes_test

import (
	"errors"
	"testing"

	"github.com/hyperboloide/pipe/encoders/aes"
	"github.com/hyperboloide/pipe/tests"
)

func TestAES(t *testing.T) {

	enc := &aes.AES{}
	if err := enc.Start(); err == nil {
		t.Error(errors.New("should validate that the key is present"))
	}

	if k, err := aes.GenKey(); err != nil {
		t.Error(err)
	} else {
		enc = &aes.AES{KeyB64: k}
	}

	if err := tests.TestEncoderDecoder(enc, "../../tests/test.jpg"); err != nil {
		t.Error(err)
	}
}
