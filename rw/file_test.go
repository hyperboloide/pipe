package rw_test

import (
	"github.com/hyperboloide/pipe/rw"
	"testing"
)

func TestFile(t *testing.T) {

	file := &rw.File{AllowSub: true}
	if err := file.Start(); err != nil {
		t.Error(err)
	}

	if err := testReadWriteDeleter(file, "some/dir/test_file"); err != nil {
		t.Error(err)
	}
}
