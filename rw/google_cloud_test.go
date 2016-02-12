package rw_test

import (
	"github.com/hyperboloide/pipe/rw"
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

	gc := &rw.GoogleCloud{
		ProjectId:   "hyperboloide",
		Bucket:      "test-pipe",
		JsonKeyPath: gcKeyPath,
	}

	if err := testReadWriteDeleter(gc, "gc_test_obj"); err != nil {
		t.Error(err)
	}
}
