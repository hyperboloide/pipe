package s3_test

import (
	"github.com/hyperboloide/pipe/rw/s3"
	"github.com/hyperboloide/pipe/tests"
	"log"
	"os"
	"testing"
)

func TestS3(t *testing.T) {

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Print("AWS_ACCESS env variable not found! Skipping TestS3")
		return
	}

	s := &s3.S3{
		Domain: "s3-eu-central-1.amazonaws.com",
		Bucket: "test.pipe",
	}

	err := tests.TestReadWriteDeleter(s, "s3_test_obj", "../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}

}
