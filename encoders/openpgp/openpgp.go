package openpgp

import (
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
)

// OpenPGP Encrypt and Decrypt with a key pair.
// The file gen.sh in the child directory 'openpgp' gives an example on
// how to generate a PGP key pair.
type OpenPGP struct {

	// If your key doesn't have a pass phrase, leave it empty.
	PassPhrase string `json:"pass_phrase"`

	// A reader to the private key file.
	// If not set then decryption will not be possible.
	PrivateKey io.Reader

	// Will read the private key from a file if set.
	PrivateKeyPath string `json:"private_key"`

	// A reader to the public key file.
	// If not set then encryption will not be possible.
	PublicKey io.Reader

	// Will read the public key from a file if set.
	PublicKeyPath string `json:"public_key"`

	privateEntityList openpgp.EntityList
	publicEntityList  openpgp.EntityList
}

// Start reads the keys and decrypt the private key if a PassPhrase is set.
func (o *OpenPGP) Start() error {
	var err error

	var privKF, pubKF *os.File

	if o.PrivateKeyPath != "" {
		if privKF, err = os.Open(o.PrivateKeyPath); err != nil {
			return err
		} else {
			o.PrivateKey = privKF
			defer privKF.Close()
		}
	}

	if o.PublicKeyPath != "" {
		if pubKF, err = os.Open(o.PublicKeyPath); err != nil {
			return err
		} else {
			o.PublicKey = pubKF
			defer pubKF.Close()
		}
	}

	if o.PrivateKey != nil {
		o.privateEntityList, err = openpgp.ReadKeyRing(o.PrivateKey)
		if err != nil {
			return err
		} else if o.PassPhrase != "" {
			ppb := []byte(o.PassPhrase)
			entity := o.privateEntityList[0]

			if err := entity.PrivateKey.Decrypt(ppb); err != nil {
				return err
			}
			for _, k := range entity.Subkeys {
				if err = k.PrivateKey.Decrypt(ppb); err != nil {
					return err
				}
			}
		}
	}

	if o.PublicKey != nil {
		o.publicEntityList, err = openpgp.ReadKeyRing(o.PublicKey)
		if err != nil {
			return err
		}
	}
	return err
}

// Encode encrypts a stream with the public key.
func (o *OpenPGP) Encode(r io.Reader, w io.Writer) error {
	if len(o.publicEntityList) == 0 {
		return errors.New("No public key defined for OpenPGP.")
	}

	wPGP, err := openpgp.Encrypt(w, o.publicEntityList, nil, nil, nil)
	if err != nil {
		return err
	}
	defer wPGP.Close()

	_, err = io.Copy(wPGP, r)
	return err
}

// Decode decrypts with the private key
func (o *OpenPGP) Decode(r io.Reader, w io.Writer) error {
	if len(o.privateEntityList) == 0 {
		return errors.New("No private key defined for OpenPGP.")
	}

	md, err := openpgp.ReadMessage(r, o.privateEntityList, nil, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, md.UnverifiedBody)
	return err
}
