package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AES encryption encrypts with a key 256 bits key using an
// aes 256 cfb scheme and an IV that is appended to the output stream.
type AES struct {
	// The encryption key must be 256 bits long.
	Key []byte

	// Alternativly to the key, a b64 encoded string can be set as the key
	// and will be decoded on start if Key is nil.
	KeyB64 string `json:"key"`
}

func (a *AES) Start() error {
	if a.Key == nil && a.KeyB64 == "" {
		return errors.New("Key undefined for AES.")
	}
	if a.Key == nil {
		if res, err := base64.StdEncoding.DecodeString(a.KeyB64); err != nil {
			return err
		} else {
			a.Key = res
		}
	}
	if len(a.Key) < 32 {
		return errors.New("Key size must be at least 32 bytes")
	}
	return nil
}

// Encode encrypts a stream with the key and generates an IV
// that will be appended to the stream.
func (a *AES) Encode(r io.Reader, w io.Writer) error {
	iv := make([]byte, aes.BlockSize)

	if block, err := aes.NewCipher(a.Key); err != nil {
		return err
	} else if _, err := rand.Read(iv[:]); err != nil {
		return err
	} else if _, err := w.Write(iv); err != nil {
		return err
	} else {
		stream := cipher.NewCFBEncrypter(block, iv)
		writer := &cipher.StreamWriter{S: stream, W: w}
		_, err := io.Copy(writer, r)
		return err
	}
}

// Decode an encrypted stream that start with an IV.
func (a *AES) Decode(r io.Reader, w io.Writer) error {
	iv := make([]byte, aes.BlockSize)

	if block, err := aes.NewCipher(a.Key); err != nil {
		return err
	} else if _, err := r.Read(iv[:]); err != nil {
		return err
	} else {
		stream := cipher.NewCFBDecrypter(block, iv)
		reader := &cipher.StreamReader{S: stream, R: r}
		_, err := io.Copy(w, reader)
		return err
	}
}

// GenKey generates a 32 bytes key encoded in base64.
func GenKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(b), nil
	}
}
