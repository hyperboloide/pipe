package gcs

import (
	"context"
	"io"

	"google.golang.org/api/option"

	"github.com/hyperboloide/pipe/rw"

	"cloud.google.com/go/storage"
)

// GCS defines a Google Cloud Storage connection to a bucket.
type GCS struct {
	rw.Prefixed

	Bucket            string `json:"bucket"`
	ServiceAccountKey string `json:"key"`
	bucket            *storage.BucketHandle
}

func (rw *GCS) Start() error {
	ctx := context.Background()
	var client *storage.Client
	var err error

	if len(rw.ServiceAccountKey) > 0 {
		opt := option.WithServiceAccountFile(rw.ServiceAccountKey)
		client, err = storage.NewClient(ctx, opt)
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return err
	}
	rw.bucket = client.Bucket(rw.Bucket)
	return nil
}

// NewWriter returns a Google Cloud Storage Writer
func (rw *GCS) NewWriter(id string) (io.WriteCloser, error) {
	ctx := context.Background()
	obj := rw.bucket.Object(rw.Prefixed.Name(id))
	return obj.NewWriter(ctx), nil
}

// NewReader returns a Google Cloud Storage Reader
func (rw *GCS) NewReader(id string) (io.ReadCloser, error) {
	ctx := context.Background()
	obj := rw.bucket.Object(rw.Prefixed.Name(id))
	return obj.NewReader(ctx)
}

// Delete an Google Cloud Storage object
func (rw *GCS) Delete(id string) error {
	ctx := context.Background()
	obj := rw.bucket.Object(rw.Prefixed.Name(id))
	return obj.Delete(ctx)
}
