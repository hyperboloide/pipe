//
// pipe.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Apr 26 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package pipe

import (
	"io"
)

type Process func(io.Reader, io.Writer) error

type Pipe struct {
	reader *io.PipeReader

	errors []chan error
	Total  int64
}

func New(reader io.Reader) *Pipe {
	r, w := io.Pipe()

	p := &Pipe{
		reader: r,
		errors: make([]chan error, 1),
	}
	p.errors[0] = make(chan error, 1)

	go func(errCh chan error) {
		total, err := io.Copy(w, reader)
		w.Close()
		p.Total = total
		errCh <- err
	}(p.errors[0])

	return p
}

func (p *Pipe) Push(procs ...Process) {
	for _, proc := range procs {
		err := make(chan error, 1)
		p.errors = append(p.errors, err)

		r, w := io.Pipe()

		go func(p Process, r io.Reader, w *io.PipeWriter, err chan error) {
			err <- p(r, w)
			w.Close()
		}(proc, p.reader, w, err)

		p.reader = r
	}
}

func (p *Pipe) To(w io.Writer) {
	go func() {
		io.Copy(w, p.reader)
		p.reader.Close()
	}()
}

func (p *Pipe) Reader() io.Reader {
	return p.reader
}

func (p *Pipe) Exec() error {
	for i, _ := range p.errors {
		if err := <-p.errors[i]; err != nil {
			close(p.errors[i])
			return err
		}
		close(p.errors[i])
	}
	return nil
}

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
