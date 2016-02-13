package encoders

import (
	"errors"
	"golang.org/x/crypto/openpgp"
	"io"
)

// Encrypt and Decrypt with an OpenPGP key pair.
// The file gen.sh in the child directory 'openpgp' gives an example on
// how to generate a PGP key pair.
type OpenPGP struct {

	// A reader to the private key file.
	// If not set then decryption will not be possible.
	PrivateKey io.Reader

	// If your key doesn't have a pass phrase, leave it empty.
	PassPhrase string

	// A reader to the public key file.
	// If not set then encryption will not be possible.
	PublicKey io.Reader

	privateEntityList openpgp.EntityList
	publicEntityList  openpgp.EntityList
}

// Reads the keys and decrypt the private key if a PassPhrase is set.
func (o *OpenPGP) Start() error {
	var err error

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

// Encrypt with the public key.
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

// Decrypt with the private key
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
