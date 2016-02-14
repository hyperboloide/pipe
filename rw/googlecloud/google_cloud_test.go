package googlecloud_test

import (
	"github.com/hyperboloide/pipe/rw/googlecloud"
	"github.com/hyperboloide/pipe/tests"
	"log"
	"os"
	"testing"
)

const (
	gcKeyPath = "./google_cloud_key.json"
)

func TestGoogleCloud(t *testing.T) {

	if _, err := os.Stat(gcKeyPath); os.IsNotExist(err) {
		log.Printf("file '%s' not found! Skipping TestGoogleCloud.", gcKeyPath)
		return
	}

	gc := &googlecloud.GoogleCloud{
		ProjectId:   "hyperboloide",
		Bucket:      "test-pipe",
		JsonKeyPath: gcKeyPath,
	}

	err := tests.TestReadWriteDeleter(gc, "gc_test_obj", "../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}
}
