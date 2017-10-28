package s3

import (
	"errors"
	"io"

	"github.com/hyperboloide/pipe/rw"
	"github.com/rlmcpherson/s3gof3r"
)

// S3DefaultDomain is the default domain to connect to S3
const S3DefaultDomain = "s3.amazonaws.com"

// S3 defines connection parameters to S3.
// An S3 Object allow the use of AWS S3.
type S3 struct {
	rw.Prefixed

	// The s3-compatible endpoint. Defaults to "s3.amazonaws.com"
	Domain string `json:"domain"`

	// If the key is not set we try to read from env
	AccessKey string `json:"access_key"` // AWS_ACCESS_KEY_ID
	SecretKey string `json:"secret_key"` // AWS_SECRET_ACCESS_KEY

	// Bucket name
	Bucket string `json:"bucket"`
	bucket *s3gof3r.Bucket
}

// Start an S3 bucket
func (s *S3) Start() error {

	if s.Domain == "" {
		s.Domain = S3DefaultDomain
	}

	if s.Bucket == "" {
		return errors.New("s3 bucket is undefined")
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

// NewWriter returns a new S3 Writer
func (s *S3) NewWriter(id string) (io.WriteCloser, error) {
	return s.bucket.PutWriter(s.Prefixed.Name(id), nil, nil)
}

// NewReader returns a new S3 Reader
func (s *S3) NewReader(id string) (io.ReadCloser, error) {
	r, _, err := s.bucket.GetReader(s.Prefixed.Name(id), nil)
	return r, err
}

// Delete an S3 object
func (s *S3) Delete(id string) error {
	return s.bucket.Delete(s.Prefixed.Name(id))
}
