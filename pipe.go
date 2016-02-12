//
// pipe.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Apr 26 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

// A simple stream processing library that works like Unix pipes.
// This library has no external dependencies and is fully asynchronous.
// Create a Pipe from a Reader, add some transformation functions and
// get the result on an io.Writer.
package pipe

import (
	"io"
)

// Pipe object
type Pipe struct {
	reader *io.PipeReader

	errors      []chan error
	errorWriter chan error

	// Total number of bytes read at the origin of the Pipe.
	TotalIn int64
	// Total number of bytes written at the end of the Pipe.
	TotalOut int64
}

// Define PipeFilter functions to transform the stream
// and add them to the Pipe.
type PipeFilter func(io.Reader, io.Writer) error

// New create a new Pipe that reads from reader.
func New(reader io.Reader) *Pipe {
	r, w := io.Pipe()

	p := &Pipe{
		reader: r,
		errors: make([]chan error, 1),
	}
	p.errors[0] = make(chan error, 1)
	p.errorWriter = make(chan error, 1)

	go func(errCh chan error) {
		total, err := io.Copy(w, reader)
		w.Close()
		p.TotalIn = total
		errCh <- err
	}(p.errors[0])

	return p
}

// Push appends a function to the Pipe.
// Note that you can add as many functions as you like at once or
// separatly. They will be processed in order.
func (p *Pipe) Push(procs ...PipeFilter) *Pipe {
	for _, proc := range procs {
		err := make(chan error, 1)
		p.errors = append(p.errors, err)

		r, w := io.Pipe()

		go func(p PipeFilter, r io.Reader, w *io.PipeWriter, err chan error) {
			err <- p(r, w)
			w.Close()
		}(proc, p.reader, w, err)

		p.reader = r
	}
	return p
}

// To writes the ouptut of the Pipe in w.
func (p *Pipe) To(w io.Writer) *Pipe {
	go func() {
		total, err := io.Copy(w, p.reader)
		p.TotalOut = total
		p.errorWriter <- err
	}()
	return p
}

// ToCloser writes the ouptut of the Pipe in io.WriteCloser w and close at the end.
func (p *Pipe) ToCloser(w io.WriteCloser) *Pipe {
	go func() {
		total, err := io.Copy(w, p.reader)
		p.TotalOut = total
		if err == nil {
			err = w.Close()
		}
		p.errorWriter <- err
	}()
	return p
}

// Exec waits for the Pipe to complete and returns an error if any
// of the functions failed.
func (p *Pipe) Exec() error {
	defer p.reader.Close()

	for i, _ := range p.errors {
		if err := <-p.errors[i]; err != nil {
			close(p.errors[i])
			return err
		}
		close(p.errors[i])
	}

	defer close(p.errorWriter)
	return <-p.errorWriter
}

// Tee creates a new Pipe to duplicate the stream.
// The stream will pass through all previously pushed functions
// before going through the tee Pipe.
// Functions pushed to the original Pipe after a call to Tee will
// not alter the new Tee Pipe.
func (p *Pipe) Tee() *Pipe {
	tR, tW := io.Pipe()

	reader := io.TeeReader(p.reader, tW)
	newR, newW := io.Pipe()

	err := make(chan error, 1)
	p.errors = append(p.errors, err)

	go func(errCh chan error) {
		_, err := io.Copy(newW, reader)
		errCh <- err
		newW.Close()
		tW.Close()
	}(err)

	newPipe := New(tR)
	p.reader = newR

	return newPipe
}
