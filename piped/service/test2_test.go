package service_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/hyperboloide/pipe/piped/service"
)

func Test2(t *testing.T) {
	var config []byte

	// create tmp dirs for testing
	destDir, err := ioutil.TempDir("", "dest")
	if err != nil {
		t.Error(err)
	} else {
		// fmt.Println(destDir)
		defer os.RemoveAll(destDir)
		cfg := fmt.Sprintf(string(fileBytes("./test2.json")[:]), destDir, destDir, destDir, destDir, destDir, destDir)
		config = []byte(cfg)
	}

	// creates the server
	r := RouterFromConfig(config, true)
	srv := httptest.NewServer(r)
	defer srv.Close()

	const id = "file_id_1234"

	// post the file
	if resp, err := http.Post(srv.URL+"/test/"+id, "image/jpeg", fileReader(testImageFile)); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 201 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if res, err := ioutil.ReadFile(destDir + "/" + id); err != nil {
		t.Error(err)
	} else if !bytes.Equal(res, fileBytes(testImageFile)) {
		t.Error(errors.New("uploaded file do not match the original"))
	}

	// get the original file
	if resp, err := http.Get(srv.URL + "/original/" + id); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if !bytes.Equal(body, fileBytes(testImageFile)) {
		t.Error(errors.New("downloaded file do not match the original"))
	}

	// get the zip file
	if resp, err := http.Get(srv.URL + "/gziped/" + id); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if !bytes.Equal(body, fileBytes(testImageFile)) {
		t.Error(errors.New("downloaded file do not match the original"))
	}

	// get the aes zip file
	if resp, err := http.Get(srv.URL + "/aes_gziped/" + id); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if !bytes.Equal(body, fileBytes(testImageFile)) {
		t.Error(errors.New("downloaded file do not match the original"))
	}

	// should not get the file if no reader declared
	if resp, err := http.Get(srv.URL + "/test/" + id); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 405 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	}

}
