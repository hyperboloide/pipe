package gzip

import (
	"compress/gzip"
	"io"
)

// Gzip is an encoders that compress a stream
type Gzip struct {
	Level int
}

// Start the Gzip encoder
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

// Encode to a Gzip stream
func (g *Gzip) Encode(r io.Reader, w io.Writer) error {
	gzw, err := gzip.NewWriterLevel(w, g.Level)
	if err != nil {
		return err
	}
	defer gzw.Close()
	_, err = io.Copy(gzw, r)
	return err
}

// Decode a Gzip stream
func (g *Gzip) Decode(r io.Reader, w io.Writer) error {
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	_, err = io.Copy(w, gzr)
	return err
}
