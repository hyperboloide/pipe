package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/encoders"
	"github.com/hyperboloide/pipe/rw"
)

// WriteOperations represents the various steps and the output necessary
// to write data.
type WriteOperations struct {
	Steps  []interface{}
	Output rw.Writer
}

// AddStep add a step to a WriteOperations from a json.RawMessage.
func (wo *WriteOperations) AddStep(js json.RawMessage) error {
	t, err := GetElementType(js)
	if err != nil {
		return err
	}
	switch t {
	case "tee":
		ts, err := NewTeeFromJSON(js)
		if err != nil {
			return err
		}
		wo.Steps = append(wo.Steps, ts)
	case "encoder":
		step, err := EncoderFromJSON(js)
		if err != nil {
			return err
		}
		wo.Steps = append(wo.Steps, step)
	case "output":
		return errors.New("element of type 'output' should appear only once as the last element of a writer")
	default:
		return fmt.Errorf("element of type '%s' is not available inside a writer", t)
	}
	return nil
}

// NewTeeFromJSON build a WriteOperations for a tee from json.
func NewTeeFromJSON(js json.RawMessage) (*WriteOperations, error) {
	tmp := struct {
		Tee json.RawMessage `json:"tee"`
	}{}
	if err := json.Unmarshal(js, &tmp); err != nil {
		return nil, err
	} else if step, err := NewWriteOperationsFromJSON(tmp.Tee); err != nil {
		return nil, err
	} else {
		return step, nil
	}
}

// NewWriteOperationsFromJSON builds a WriteOperations from json.
func NewWriteOperationsFromJSON(js json.RawMessage) (*WriteOperations, error) {
	ops := []json.RawMessage{}
	if err := json.Unmarshal(js, &ops); err != nil {
		return nil, err
	}
	if len(ops) < 1 {
		return nil, errors.New("writer cannot be empty, it should end with an output")
	}
	res := &WriteOperations{}
	for i := 0; i < (len(ops) - 1); i++ {
		if err := res.AddStep(ops[i]); err != nil {
			return nil, err
		}
	}
	var w rw.Writer
	if t, err := GetElementType(ops[len(ops)-1]); err != nil {
		return nil, err
	} else if t != "output" {
		return nil, errors.New("a writer should end with an output")
	} else if w, err = WriterFromJSON(ops[len(ops)-1]); err != nil {
		return nil, err
	}
	res.Output = w
	return res, nil
}

// SetPipe set the various encoders, tees and the writer.
func (wo *WriteOperations) SetPipe(p *pipe.Pipe, id string) error {
	for _, s := range wo.Steps {
		switch s.(type) {
		case encoders.Encoder:
			p.Push(s.(encoders.Encoder).Encode)
		case *WriteOperations:
			tp := p.Tee()
			teeWo := s.(*WriteOperations)
			if err := teeWo.SetPipe(tp, id); err != nil {
				return err
			}
		}
	}
	w, err := wo.Output.NewWriter(id)
	if err != nil {
		return err
	}
	p.ToCloser(w)
	return nil
}

// ReadOperations represents the various steps and the input necessary
// to retrieve data.
type ReadOperations struct {
	Steps []encoders.Decoder
	Input rw.Reader
}

// AddStep add a step to a ReadOperations from a json.RawMessage.
func (ro *ReadOperations) AddStep(js json.RawMessage) error {
	t, err := GetElementType(js)
	if err != nil {
		return err
	}
	switch t {
	case "decoder":
		step, err := DecoderFromJSON(js)
		if err != nil {
			return err
		}
		ro.Steps = append(ro.Steps, step)
	case "input":
		return fmt.Errorf("element of type '%s' is not available inside a reader", t)
	default:
		return errors.New("element of type 'input' should appear only once as the first element of a reader")
	}
	return nil
}

// NewReadOperationsFromJSON builds a ReadOperations from json.
func NewReadOperationsFromJSON(js json.RawMessage) (*ReadOperations, error) {
	ops := []json.RawMessage{}
	if err := json.Unmarshal(js, &ops); err != nil {
		return nil, err
	}
	if len(ops) < 1 {
		return nil, errors.New("reader cannot be empty, it should start with an input")
	}

	res := &ReadOperations{}
	var r rw.Reader
	if t, err := GetElementType(ops[0]); err != nil {
		return nil, err
	} else if t != "input" {
		return nil, errors.New("a reader should start with an input")
	} else if r, err = ReaderFromJSON(ops[0]); err != nil {
		return nil, err
	}
	res.Input = r
	for i := 1; i < len(ops); i++ {
		if err := res.AddStep(ops[i]); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// SetPipe adds the decoders to the pipe.
func (ro *ReadOperations) SetPipe(p *pipe.Pipe) error {
	for _, s := range ro.Steps {
		p.Push(s.Decode)
	}
	return nil
}
