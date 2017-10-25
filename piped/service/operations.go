package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/encoders"
	"github.com/hyperboloide/pipe/rw"
)

type WriteOperations struct {
	Steps  []interface{}
	Output rw.Writer
}

func NewWriteOperationsFromJson(js json.RawMessage) (*WriteOperations, error) {
	ops := []json.RawMessage{}
	if err := json.Unmarshal(js, &ops); err != nil {
		return nil, err
	}
	if len(ops) < 1 {
		return nil, errors.New("Writer cannot be empty, it should end with an output.")
	}
	res := &WriteOperations{}
	for i := 0; i < (len(ops) - 1); i++ {

		if t, err := GetElementType(ops[i]); err != nil {
			return nil, err
		} else if t == "tee" {
			tmp := struct {
				Tee json.RawMessage `json:"tee"`
			}{}
			if err := json.Unmarshal(ops[i], &tmp); err != nil {
				return nil, err
			} else if step, err := NewWriteOperationsFromJson(tmp.Tee); err != nil {
				return nil, err
			} else {
				res.Steps = append(res.Steps, step)
			}
		} else if t == "encoder" {
			if step, err := EncoderFromJson(ops[i]); err != nil {
				return nil, err
			} else {
				res.Steps = append(res.Steps, step)
			}
		} else if t == "output" {
			return nil, errors.New("element of type 'output' should appear only once as the last element of a writer.")
		} else {
			return nil, fmt.Errorf("element of type '%s' is not available inside a writer", t)
		}
	}
	if t, err := GetElementType(ops[len(ops)-1]); err != nil {
		return nil, err
	} else if t != "output" {
		return nil, errors.New("a writer should end with an output.")
	} else if w, err := WriterFromJson(ops[len(ops)-1]); err != nil {
		return nil, err
	} else {
		res.Output = w
		return res, nil
	}
}

func (wo *WriteOperations) SetPipe(p *pipe.Pipe, id string) error {
	for _, s := range wo.Steps {
		switch s.(type) {
		case encoders.Encoder:
			p.Push(s.(encoders.Encoder).Encode)
		case *WriteOperations:
			tp := p.Tee()
			wo := s.(*WriteOperations)
			wo.SetPipe(tp, id)
		}
	}
	if w, err := wo.Output.NewWriter(id); err != nil {
		return err
	} else {
		p.ToCloser(w)
		return nil
	}
}

type ReadOperations struct {
	Steps []encoders.Decoder
	Input rw.Reader
}

func NewReadOperationsFromJson(js json.RawMessage) (*ReadOperations, error) {
	ops := []json.RawMessage{}
	if err := json.Unmarshal(js, &ops); err != nil {
		return nil, err
	}
	if len(ops) < 1 {
		return nil, errors.New("Reader cannot be empty, it should start with an input.")
	}

	res := &ReadOperations{}

	if t, err := GetElementType(ops[0]); err != nil {
		return nil, err
	} else if t != "input" {
		return nil, errors.New("a reader should start with an input.")
	} else if r, err := ReaderFromJson(ops[0]); err != nil {
		return nil, err
	} else {
		res.Input = r
	}

	for i := 1; i < len(ops); i++ {
		if t, err := GetElementType(ops[i]); err != nil {
			return nil, err
		} else if t == "decoder" {
			if step, err := DecoderFromJson(ops[i]); err != nil {
				return nil, err
			} else {
				res.Steps = append(res.Steps, step)
			}
		} else if t == "input" {
			return nil, errors.New("element of type 'input' should appear only once as the first element of a reader.")
		} else {
			return nil, fmt.Errorf("element of type '%s' is not available inside a reader", t)
		}

	}

	return res, nil
}

func (ro *ReadOperations) SetPipe(p *pipe.Pipe) error {
	for _, s := range ro.Steps {
		p.Push(s.Decode)
	}
	return nil
}
