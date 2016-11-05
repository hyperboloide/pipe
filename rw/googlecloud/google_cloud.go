package googlecloud

import (
	"errors"
	"github.com/hyperboloide/pipe/rw"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"io"
	"io/ioutil"
)

// GoogleCloudScope is the scope to access the bucket.
const GoogleCloudScope = "https://www.googleapis.com/auth/devstorage.read_write"

var (
	// ErrGcBucketUndefined is returned if the bucket is undefined.
	ErrGcBucketUndefined = errors.New("Google Cloud bucket is undefined.")

	// ErrGcProjectIDUndefiend is returned if the
	// Google Cloud project id is undefined.
	ErrGcProjectIDUndefiend = errors.New("Google Cloud project id is undefined.")

	// ErrNoAuthProvided is returned if no mean of identification is provided.
	ErrNoAuthProvided = errors.New("No mean of identification provided.")

	// ErrJSONKeyNotFound is returned if the json key cannot be read
	ErrJSONKeyNotFound = errors.New("Cannot read json key.")
)

// GoogleCloud defines connection parameters to Google Cloud Storage
type GoogleCloud struct {
	rw.Prefixed
	ProjectID    string
	Bucket       string
	JSONKeyPath  string
	context      context.Context
	bucketHandle *storage.BucketHandle
}

// Start a connection to Google Cloud Storage
func (g *GoogleCloud) Start() error {
	if g.Bucket == "" {
		return ErrGcBucketUndefined
	} else if g.ProjectID == "" {
		return ErrGcProjectIDUndefiend
	}

	if g.JSONKeyPath == "" {
		return ErrNoAuthProvided
	}
	data, err := ioutil.ReadFile(g.JSONKeyPath)
	if err != nil {
		return ErrJSONKeyNotFound
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

// NewWriter returns a Writer to Google Cloud Storage
func (g *GoogleCloud) NewWriter(id string) (io.WriteCloser, error) {
	return g.bucketHandle.Object(g.Prefixed.Name(id)).NewWriter(nil), nil
}

// NewReader returns a reader to Google Cloud Storage
func (g *GoogleCloud) NewReader(id string) (io.ReadCloser, error) {
	return g.bucketHandle.Object(g.Prefixed.Name(id)).NewReader(nil)
}

// Delete an object on Google Cloud Storage
func (g *GoogleCloud) Delete(id string) error {
	return g.bucketHandle.Object(g.Prefixed.Name(id)).Delete(nil)
}
