package s3_test

import (
	"github.com/hyperboloide/pipe/rw/s3"
	"github.com/hyperboloide/pipe/tests"
	"log"
	"os"
	"testing"
)

func TestS3(t *testing.T) {

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		log.Fatal("AWS_ACCESS_KEY_ID  env variable not found! Skipping TestS3")

	} else if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatal("AWS_SECRET_ACCESS_KEY  env variable not found! Skipping TestS3")
	}

	s := &s3.S3{
		Bucket:    "test.pipe",
		Domain:    "s3-eu-central-1.amazonaws.com",
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if err := s.Start(); err != nil {
		t.Error(err)
	}

	err := tests.TestReadWriteDeleter(s, "s3_test_obj", "../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}

}
