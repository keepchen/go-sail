package utils

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	contents = []byte(`hello
world
hello
go-sail`)
	dst = "./file_test_will_be_deleted_after_testcase.txt"
)

func createFormFile() *multipart.FileHeader {
	// Create a buffer to store multipart data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add a file field (e.g. named "file")
	part, err := writer.CreateFormFile("file", "example.txt")
	if err != nil {
		panic(err)
	}

	// Write file content
	content := []byte("Hello, World!")
	_, _ = part.Write(content)

	// Add other form fields (optional)
	_ = writer.WriteField("key", "value")

	// Close the writer to finish writing data
	_ = writer.Close()

	// Create an HTTP request and parse multipart data
	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Parse multipart data in the request
	_ = req.ParseMultipartForm(32 << 20) // 32MB memory limit

	return req.MultipartForm.File["file"][0]
}

func TestSave2Dst(t *testing.T) {
	t.Run("Save2Dst", func(t *testing.T) {
		sDst := "./temp-file"
		f := &multipart.FileHeader{}
		t.Log(File().Save2Dst(f, sDst))
	})

	t.Run("Save2Dst-Copy", func(t *testing.T) {
		f := createFormFile()
		t.Log(File().Save2Dst(f, dst))
		cleanClear(nil)
	})
}

func TestFileImplGetContents(t *testing.T) {
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	_, err2 := File().GetContents(dst)
	assert.Equal(t, nil, err2)
	cleanClear(err)
}

func TestFileImplPutContents(t *testing.T) {
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	cleanClear(err)
}

func TestFileImplAppendContents(t *testing.T) {
	t.Run("AppendContents", func(t *testing.T) {
		err := File().PutContents(contents, dst)
		assert.Equal(t, nil, err)
		err2 := File().AppendContents([]byte(`new line`), dst)
		assert.Equal(t, nil, err2)
		cleanClear(err)
	})

	t.Run("AppendContents-Error", func(t *testing.T) {
		err := File().AppendContents([]byte(`new line`), dst)
		assert.Equal(t, nil, err)
		cleanClear(err)
	})
}

func TestFileImplExists(t *testing.T) {
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, File().Exists(dst))
	cleanClear(err)
}

func TestFileImplExistsWithError(t *testing.T) {
	ok, err := File().ExistsWithError(dst)
	assert.Equal(t, false, ok)
	cleanClear(err)
}

func TestFileImplExt(t *testing.T) {
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	assert.Equal(t, "txt", File().Ext(dst))
	cleanClear(err)
}

func TestFileImplGetContentsReadLine(t *testing.T) {
	t.Run("GetContentsReadLine", func(t *testing.T) {
		err := File().PutContents(contents, dst)
		assert.Equal(t, nil, err)

		ch, err2 := File().GetContentsReadLine(dst)
		assert.Equal(t, nil, err2)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for content := range ch {
				log.Printf("[TestFileImplGetContentsReadLine] content: %s\n", content)
				assert.NotEqual(t, "", content)
			}
		}()
		wg.Wait()
		cleanClear(err)
	})

	t.Run("GetContentsReadLine-Panic", func(t *testing.T) {
		ch, err := File().GetContentsReadLine(dst)
		assert.Error(t, err)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for content := range ch {
				log.Printf("[TestFileImplGetContentsReadLine] content: %s\n", content)
				assert.NotEqual(t, "", content)
			}
		}()
		wg.Wait()
		cleanClear(err)
	})
}

func cleanClear(err error) {
	if err == nil {
		_ = os.RemoveAll(dst)
	}
}
