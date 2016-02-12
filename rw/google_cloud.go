package rw

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"io"
	"io/ioutil"
)

const GoogleCloudScope = "https://www.googleapis.com/auth/devstorage.read_write"

var (
	GcBucketUndefined    = errors.New("Google Cloud bucket is undefined.")
	GcProjectIdUndefiend = errors.New("Google Cloud project id is undefined.")
	NoAuthProvided       = errors.New("No mean of identification provided.")
	JsonKeyNotFound      = errors.New("Cannot read json key.")
)

// Defines connection parameters to Google Cloud Storage
type GoogleCloud struct {
	Prefixed
	ProjectId    string
	Bucket       string
	JsonKeyPath  string
	context      context.Context
	bucketHandle *storage.BucketHandle
}

// Starts a connection to Google Cloud Storage
func (g *GoogleCloud) Start() error {
	if g.Bucket == "" {
		return GcBucketUndefined
	} else if g.ProjectId == "" {
		return GcProjectIdUndefiend
	}

	if g.JsonKeyPath == "" {
		return NoAuthProvided
	}
	data, err := ioutil.ReadFile(g.JsonKeyPath)
	if err != nil {
		return JsonKeyNotFound
	}
	conf, err := google.JWTConfigFromJSON(data, GoogleCloudScope)
	if err != nil {
		return err
	}

	g.context = context.Background()

	client, err := storage.NewClient(g.context, cloud.WithTokenSource(conf.TokenSource(g.context)))
	if err != nil {
		return err
	}
	g.bucketHandle = client.Bucket(g.Bucket)
	return nil
}

// Returns a Writer to Google Cloud Storage
func (g *GoogleCloud) NewWriter(id string) (io.WriteCloser, error) {
	oh := g.bucketHandle.Object(g.Prefixed.Name(id))
	return oh.NewWriter(nil), nil
	// return storage.NewWriter(g.context, g.Bucket, g.Prefixed.Name(id)), nil
}

// Returns a reader to Google Cloud Storage
func (g *GoogleCloud) NewReader(id string) (io.ReadCloser, error) {
	oh := g.bucketHandle.Object(g.Prefixed.Name(id))
	return oh.NewReader(nil)
}

// Deletes and object on Google Cloud Storage
func (g *GoogleCloud) Delete(id string) error {
	oh := g.bucketHandle.Object(g.Prefixed.Name(id))
	return oh.Delete(nil)
}
