package encoders

import (
	"compress/gzip"
	"io"
)

type Gzip struct {
	Level int
}

func (g *Gzip) Start() error {
	switch g.Level {
	case gzip.DefaultCompression:
	case gzip.BestCompression:
	case gzip.BestSpeed:
	default:
		g.Level = gzip.DefaultCompression
	}
	return nil
}

func (g *Gzip) Encode(r io.Reader, w io.Writer) error {
	gzw, err := gzip.NewWriterLevel(w, g.Level)
	if err != nil {
		return err
	}
	defer gzw.Close()
	_, err = io.Copy(gzw, r)
	return err
}

func (g *Gzip) Decode(r io.Reader, w io.Writer) error {
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	_, err = io.Copy(w, gzr)
	return err
}
