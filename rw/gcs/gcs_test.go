package gcs_test

import (
	"testing"

	"github.com/hyperboloide/pipe/rw/gcs"
	"github.com/hyperboloide/pipe/tests"
)

func TestGCS(t *testing.T) {

	s := &gcs.GCS{
		Bucket:            "hyperboloide-pipe-test",
		ServiceAccountKey: "./Ozigo-aab90df49241.json",
	}

	if err := s.Start(); err != nil {
		t.Error(err)
	}

	err := tests.TestReadWriteDeleter(s, "pipe_test_obj", "../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}

}
