package s3

import (
	"errors"
	"github.com/hyperboloide/pipe/rw"
	"github.com/rlmcpherson/s3gof3r"
	"io"
)

const S3DefaultDomain = "s3.amazonaws.com"

var (
	S3BucketUndefined = errors.New("S3 bucket is undefined.")
)

// Defines connection parameters to S3.
// An S3 Object allow the use of AWS S3.
type S3 struct {
	rw.Prefixed

	// The s3-compatible endpoint. Defaults to "s3.amazonaws.com"
	Domain string

	// If the key is not set we try to read from env
	AccessKey string // AWS_ACCESS_KEY_ID
	SecretKey string // AWS_SECRET_ACCESS_KEY

	// Bucket name
	Bucket string
	bucket *s3gof3r.Bucket
}

// Starts an S3 bucket
func (s *S3) Start() error {

	if s.Domain == "" {
		s.Domain = S3DefaultDomain
	}

	if s.Bucket == "" {
		return S3BucketUndefined
	}

	var s3p *s3gof3r.S3
	if s.AccessKey != "" && s.SecretKey != "" {
		s3p = s3gof3r.New(
			s.Domain,
			s3gof3r.Keys{
				AccessKey: s.AccessKey,
				SecretKey: s.SecretKey})
	} else if keys, err := s3gof3r.EnvKeys(); err != nil {
		return err
	} else {
		s3p = s3gof3r.New(s.Domain, keys)
	}
	s.bucket = s3p.Bucket(s.Bucket)
	return nil
}

// Returns a new S3 Writer
func (s *S3) NewWriter(id string) (io.WriteCloser, error) {
	return s.bucket.PutWriter(s.Prefixed.Name(id), nil, nil)
}

// Returns a new S3 Reader
func (s *S3) NewReader(id string) (io.ReadCloser, error) {
	r, _, err := s.bucket.GetReader(s.Prefixed.Name(id), nil)
	return r, err
}

// Delete an S3 object
func (s *S3) Delete(id string) error {
	return s.bucket.Delete(s.Prefixed.Name(id))
}
