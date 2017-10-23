package gcs_test

import (
	"os"
	"testing"

	"github.com/hyperboloide/pipe/rw/gcs"
	"github.com/hyperboloide/pipe/tests"
)

func TestGCS(t *testing.T) {

	if os.Getenv("GCS_BUCKET") == "" {
		t.Skip("GCS_BUCKET env variable not set! Skipping Test")
	} else if os.Getenv("GCS_KEY_FILE") == "" {
		t.Skip("GCS_KEY_FILE env variable not set! Skipping Test")
	}

	s := &gcs.GCS{
		Bucket:            os.Getenv("GCS_BUCKET"),
		ServiceAccountKey: os.Getenv("GCS_KEY_FILE"),
	}

	if err := s.Start(); err != nil {
		t.Error(err)
	}

	err := tests.TestReadWriteDeleter(s, "pipe_test_obj", "../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}

}
