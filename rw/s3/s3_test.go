package s3_test

import (
	"os"
	"testing"

	"github.com/hyperboloide/pipe/rw/s3"
	"github.com/hyperboloide/pipe/tests"
)

func TestS3(t *testing.T) {

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("AWS_ACCESS_KEY_ID env variable not set! Skipping Test")
	} else if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.Skip("AWS_SECRET_ACCESS_KEY env variable not set! Skipping Test")
	} else if os.Getenv("AWS_S3_BUCKET") == "" {
		t.Skip("AWS_S3_BUCKET env variable not set! Skipping Test")
	} else if os.Getenv("AWS_S3_DOMAIN") == "" {
		t.Skip("AWS_S3_DOMAIN env variable not set! Skipping Test")
	}

	s := &s3.S3{
		Bucket:    os.Getenv("AWS_S3_BUCKET"),
		Domain:    os.Getenv("AWS_S3_DOMAIN"),
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
