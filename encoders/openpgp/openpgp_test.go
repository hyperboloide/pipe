package openpgp_test

import (
	"testing"

	"github.com/hyperboloide/pipe/encoders/openpgp"
	"github.com/hyperboloide/pipe/tests"
)

func TestOpenPGP(t *testing.T) {

	enc := &openpgp.OpenPGP{
		PrivateKeyPath: "./the_key.sec",
		PublicKeyPath:  "./the_key.pub",
	}

	if err := tests.TestEncoderDecoder(enc, "../../tests/test.jpg"); err != nil {
		t.Error(err)
	}
}
