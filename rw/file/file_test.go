package file_test

import (
	"github.com/hyperboloide/pipe/rw/file"
	"github.com/hyperboloide/pipe/tests"
	"testing"
)

func TestFile(t *testing.T) {

	file := &file.File{AllowSub: true}
	if err := file.Start(); err != nil {
		t.Error(err)
	}

	err := tests.TestReadWriteDeleter(file, "some/dir/test_file", "../../tests/test.jpg")
	if err != nil {
		t.Error(err)
	}
}
