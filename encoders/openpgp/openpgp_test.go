package openpgp_test

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/encoders/openpgp"
	"io/ioutil"
	"os"
	"testing"
)

func TestOpenPGP(t *testing.T) {

	rPub, err := os.Open("./the_key.pub")
	if err != nil {
		t.Error(err)
	}

	rPriv, err := os.Open("./the_key.sec")
	if err != nil {
		t.Error(err)
	}

	pgp := &openpgp.OpenPGP{
		PrivateKey: rPriv,
		PublicKey:  rPub,
	}
	if err := pgp.Start(); err != nil {
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

	if err := pipe.New(reader).Push(pgp.Encode).To(writer).Exec(); err != nil {
		t.Error(err)
	} else if bytes.Equal(result1.Bytes(), bin) == true {
		t.Error(errors.New("Content match original."))
	}

	reader2 := bytes.NewReader(result1.Bytes())
	var result2 bytes.Buffer
	writer = bufio.NewWriter(&result2)

	if err := pipe.New(reader2).Push(pgp.Decode).To(writer).Exec(); err != nil {
		t.Error(err)
	} else if bytes.Equal(result1.Bytes(), bin) == true {
		t.Error(errors.New("Content match original."))
	}
}
