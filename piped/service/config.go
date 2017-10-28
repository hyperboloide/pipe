package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperboloide/pipe/encoders"
	"github.com/hyperboloide/pipe/encoders/aes"
	"github.com/hyperboloide/pipe/encoders/gzip"
	"github.com/hyperboloide/pipe/encoders/openpgp"
	"github.com/hyperboloide/pipe/rw"
	"github.com/hyperboloide/pipe/rw/file"
	"github.com/hyperboloide/pipe/rw/gcs"
	"github.com/hyperboloide/pipe/rw/s3"
)

var (
	// ErrUnknowElementType is returned if an element does not implement a valid type.
	ErrUnknowElementType = errors.New("unknow element does not implement a valid type")
)

// GetElementType return a string representing the element type.
func GetElementType(js json.RawMessage) (string, error) {
	tmp := map[string]interface{}{}
	types := []string{"tee", "encoder", "decoder", "input", "output"}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return "", err
	}
	for _, t := range types {
		if _, ok := tmp[t]; ok {
			return t, nil
		}
	}
	return "", ErrUnknowElementType
}

// Startable is an object that implements a Start() method. It'is used
// to start encoders and rws.
type Startable interface {
	Start() error
}

// UnmarshalAndStart unmarshals a Startable element and start it.
func UnmarshalAndStart(s Startable, js json.RawMessage) error {
	if err := json.Unmarshal(js, s); err != nil {
		return err
	} else if err := s.Start(); err != nil {
		return err
	}
	return nil
}

// EncoderFromJSON builds an encoder from json.
func EncoderFromJSON(js json.RawMessage) (encoders.Encoder, error) {
	tmp := struct {
		Type string `json:"encoder"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else if res := EncoderDecoderFromString(tmp.Type); res == nil {
		return nil, fmt.Errorf("encoder of type '%s' is not supported", tmp.Type)
	} else {
		return res, UnmarshalAndStart(res, js)
	}
}

// DecoderFromJSON builds a decoder from json.
func DecoderFromJSON(js json.RawMessage) (encoders.Decoder, error) {
	tmp := struct {
		Type string `json:"decoder"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else if res := EncoderDecoderFromString(tmp.Type); res == nil {
		return nil, fmt.Errorf("decoder of type '%s' is not supported", tmp.Type)
	} else {
		return res, UnmarshalAndStart(res, js)
	}
}

// WriterFromJSON builds a writer from json.
func WriterFromJSON(js json.RawMessage) (rw.Writer, error) {
	tmp := struct {
		Type string `json:"output"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else if res := RWDFromString(tmp.Type); res == nil {
		return nil, fmt.Errorf("output of type '%s' is not supported", tmp.Type)
	} else {
		return res, UnmarshalAndStart(res, js)
	}
}

// ReaderFromJSON builds a reader from json.
func ReaderFromJSON(js json.RawMessage) (rw.Reader, error) {
	tmp := struct {
		Type string `json:"input"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else if res := RWDFromString(tmp.Type); res == nil {
		return nil, fmt.Errorf("input of type '%s' is not supported", tmp.Type)
	} else {
		return res, UnmarshalAndStart(res, js)
	}
}

// DeleterFromJSON builds a deleter from json.
func DeleterFromJSON(js json.RawMessage) (rw.Deleter, error) {
	tmp := struct {
		Type string `json:"type"`
	}{}
	var res rw.Deleter
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else if tmp.Type == "" {
		return nil, errors.New("a deleter should define it's type")
	} else if res = RWDFromString(tmp.Type); res == nil {
		return nil, fmt.Errorf("deleter of type '%s' is not supported", tmp.Type)
	}
	return res, UnmarshalAndStart(res, js)
}

// RWDFromString returns a rw.ReadWriteDeleter from it's name.
func RWDFromString(str string) rw.ReadWriteDeleter {
	var res rw.ReadWriteDeleter
	switch str {
	case "file":
		res = &file.File{}
	case "gcs":
		res = &gcs.GCS{}
	case "s3":
		res = &s3.S3{}
	}
	return res
}

// EncoderDecoderFromString returns an encoders.EncoderDecoder from it's name.
func EncoderDecoderFromString(str string) encoders.EncoderDecoder {
	var res encoders.EncoderDecoder
	switch str {
	case "gzip":
		res = &gzip.Gzip{}
	case "aes":
		res = &aes.AES{}
	case "openpgp":
		res = &openpgp.OpenPGP{}
	}
	return res
}
