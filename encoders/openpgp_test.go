package encoders_test

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/encoders"
	"log"
	"os"
	"testing"
)

func TestOpenPGP(t *testing.T) {

	rPub, err := os.Open("./openpgp/the_key.pub")
	if err != nil {
		t.Error(err)
	}

	rPriv, err := os.Open("./openpgp/the_key.sec")
	if err != nil {
		t.Error(err)
	}

	pgp := &encoders.OpenPGP{
		PrivateKey: rPriv,
		PublicKey:  rPub,
	}
	if err := pgp.Start(); err != nil {
		t.Error(err)
	}

	reader := bytes.NewReader(bin)
	var result1 bytes.Buffer
	writer := bufio.NewWriter(&result1)

	if err := pipe.New(reader).Push(pgp.Encode).To(writer).Exec(); err != nil {
		t.Error(err)
	} else if bytes.Equal(result1.Bytes(), bin) == true {
		t.Error(errors.New("Content match original."))
	}

	log.Print(string(bin[:]))

	reader = bytes.NewReader(result1.Bytes())
	var result2 bytes.Buffer
	writer = bufio.NewWriter(&result2)

	if err := pipe.New(reader).Push(pgp.Decode).To(writer).Exec(); err != nil {
		t.Error(err)
	} else if bytes.Equal(result1.Bytes(), bin) == true {
		t.Error(errors.New("Content match original."))
	}
}
