package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/hyperboloide/pipe/piped/service"
)

const testImageFile = "../../tests/test.jpg"
const testTextFile = "../../tests/test.txt"

func fileReader(pth string) *os.File {
	f, err := os.Open(pth)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func fileBytes(pth string) []byte {
	b, err := ioutil.ReadFile(pth)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func Test1(t *testing.T) {
	var config []byte

	// create tmp dirs for testing
	destDir, err := ioutil.TempDir("", "dest")
	if err != nil {
		t.Error(err)
	} else {
		defer os.RemoveAll(destDir)
		cfg := fmt.Sprintf(string(fileBytes("./test1.json")[:]), destDir, destDir, destDir)
		config = []byte(cfg)
	}

	// creates the server
	r := RouterFromConfig(config, true)
	srv := httptest.NewServer(r)
	defer srv.Close()

	const id = "file_id_1234"
	data := &WriteResponse{}

	// post the file
	if resp, err := http.Post(srv.URL+"/test/"+id, "image/jpeg", fileReader(testImageFile)); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 201 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if res, err := ioutil.ReadFile(destDir + "/" + id); err != nil {
		t.Error(err)
	} else if !bytes.Equal(res, fileBytes(testImageFile)) {
		t.Error(errors.New("uploaded file do not match the original"))
	} else if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Error(err)
	} else if int64(len(fileBytes(testImageFile))) != data.BytesIn {
		t.Error(errors.New("result size of file do not match the original"))
	}

	// get the file
	if resp, err := http.Get(srv.URL + "/test/" + id); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else if !bytes.Equal(body, fileBytes(testImageFile)) {
		t.Error(errors.New("downloaded file do not match the original"))
	}

	// update the file
	if resp, err := http.Post(srv.URL+"/test/"+id, "text/plain", fileReader(testTextFile)); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 201 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if res, err := ioutil.ReadFile(destDir + "/" + id); err != nil {
		t.Error(err)
	} else if !bytes.Equal(res, fileBytes(testTextFile)) {
		t.Error(errors.New("updated file do not match the original"))
	}

	// delete the file
	client := &http.Client{}
	if req, err := http.NewRequest("DELETE", srv.URL+"/test/"+id, nil); err != nil {
		t.Error(err)
	} else if resp, err := client.Do(req); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 204 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if _, err := os.Stat(destDir + "/" + id); !os.IsNotExist(err) {
		t.Error(errors.New("file not deleted"))
	}

	// get that dont exists
	if resp, err := http.Get(srv.URL + "/test/" + id); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 404 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	}

	// delete the file that dont exists
	client = &http.Client{}
	if req, err := http.NewRequest("DELETE", srv.URL+"/test/"+id, nil); err != nil {
		t.Error(err)
	} else if resp, err := client.Do(req); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 204 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	}

	// post a file with no id
	if files, err := ioutil.ReadDir(destDir); err != nil {
		t.Error(err)
	} else if len(files) != 0 {
		t.Error(errors.New("directory should be empty"))
	} else if resp, err := http.Post(srv.URL+"/test", "image/jpeg", fileReader(testImageFile)); err != nil {
		t.Error(err)
	} else if resp.StatusCode != 201 {
		t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
	} else if files, err := ioutil.ReadDir(destDir); err != nil {
		t.Error(err)
	} else if len(files) != 1 {
		t.Error(errors.New("directory should have 1 file"))
	} else if !bytes.Equal(fileBytes(testImageFile), fileBytes(destDir+"/"+files[0].Name())) {
		t.Error(errors.New("updated file do not match the original"))
	}
}

// test form upload
func Test1Form(t *testing.T) {
	var config []byte

	// create tmp dirs for testing
	destDir, err := ioutil.TempDir("", "dest")
	if err != nil {
		t.Error(err)
	} else {
		defer os.RemoveAll(destDir)
		cfg := fmt.Sprintf(string(fileBytes("./test1.json")[:]), destDir, destDir, destDir)
		config = []byte(cfg)
	}

	// creates the server
	r := RouterFromConfig(config, true)
	srv := httptest.NewServer(r)
	defer srv.Close()

	const id = "file_from_form"

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	if fileWriter, err := bodyWriter.CreateFormFile("file", "image.jpg"); err != nil {
		t.Error(err)
	} else if fh, err := os.Open(testImageFile); err != nil {
		t.Error(err)
	} else if _, err = io.Copy(fileWriter, fh); err != nil {
		t.Error(err)
	} else {
		contentType := bodyWriter.FormDataContentType()
		bodyWriter.Close()

		if resp, err := http.Post(srv.URL+"/test/"+id, contentType, bodyBuf); err != nil {
			t.Error(err)
		} else if resp.StatusCode != 201 {
			t.Error(fmt.Errorf("invalid response status code %d", resp.StatusCode))
		} else if res, err := ioutil.ReadFile(destDir + "/" + id); err != nil {
			t.Error(err)
		} else if !bytes.Equal(res, fileBytes(testImageFile)) {
			t.Error(errors.New("uploaded file do not match the original"))
		}
	}

}
