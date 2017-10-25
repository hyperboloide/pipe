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
	ErrUnknowElementType = errors.New("unknow element does not implement a valid type.")
)

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

func GetType(js json.RawMessage) (string, error) {
	tmp := struct {
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(js, &tmp)
	return tmp.Type, err
}

type Startable interface {
	Start() error
}

func UnmarshalAndStart(s Startable, js json.RawMessage) error {
	if err := json.Unmarshal(js, s); err != nil {
		return err
	} else if err := s.Start(); err != nil {
		return err
	}
	return nil
}

func EncoderFromJson(js json.RawMessage) (encoders.Encoder, error) {
	tmp := struct {
		Type string `json:"encoder"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else {
		var res encoders.Encoder
		switch tmp.Type {
		case "gzip":
			res = &gzip.Gzip{}
		case "aes":
			res = &aes.AES{}
		case "openpgp":
			res = &openpgp.OpenPGP{}
		default:
			return nil, fmt.Errorf("Encoder of type '%s' is not supported", tmp.Type)
		}
		return res, UnmarshalAndStart(res, js)
	}
}

func DecoderFromJson(js json.RawMessage) (encoders.Decoder, error) {
	tmp := struct {
		Type string `json:"decoder"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else {
		var res encoders.Decoder
		switch tmp.Type {
		case "gzip":
			res = &gzip.Gzip{}
		case "aes":
			res = &aes.AES{}
		case "openpgp":
			res = &openpgp.OpenPGP{}
		default:
			return nil, fmt.Errorf("Decoder of type '%s' is not supported", tmp.Type)
		}
		return res, UnmarshalAndStart(res, js)
	}
}

func WriterFromJson(js json.RawMessage) (rw.Writer, error) {
	tmp := struct {
		Type string `json:"output"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else {
		var res rw.Writer
		switch tmp.Type {
		case "file":
			res = &file.File{}
		case "gcs":
			res = &gcs.GCS{}
		case "s3":
			res = &s3.S3{}
		default:
			return nil, fmt.Errorf("Output of type '%s' is not supported", tmp.Type)
		}
		return res, UnmarshalAndStart(res, js)
	}
}

func ReaderFromJson(js json.RawMessage) (rw.Reader, error) {
	tmp := struct {
		Type string `json:"input"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else {
		var res rw.Reader
		switch tmp.Type {
		case "file":
			res = &file.File{}
		case "gcs":
			res = &gcs.GCS{}
		case "s3":
			res = &s3.S3{}
		default:
			return nil, fmt.Errorf("Input of type '%s' is not supported", tmp.Type)
		}
		return res, UnmarshalAndStart(res, js)
	}
}

func DeleterFromJson(js json.RawMessage) (rw.Deleter, error) {
	tmp := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else {
		var res rw.Deleter
		switch tmp.Type {
		case "file":
			res = &file.File{}
		case "gcs":
			res = &gcs.GCS{}
		case "s3":
			res = &s3.S3{}
		case "":
			return nil, errors.New("A deleter should define it's type.")
		default:
			return nil, fmt.Errorf("Deleter of type '%s' is not supported", tmp.Type)
		}
		return res, UnmarshalAndStart(res, js)
	}
}
