package rw_test

import (
	"github.com/hyperboloide/pipe/rw"
	"log"
	"os"
	"testing"
)

func TestS3(t *testing.T) {

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Print("AWS_ACCESS env variable not found! Skipping TestS3")
		return
	}

	s3 := &rw.S3{
		Domain: "s3-eu-central-1.amazonaws.com",
		Bucket: "test.pipe",
	}

	if err := testReadWriteDeleter(s3, "s3_test_obj"); err != nil {
		t.Error(err)
	}

}
