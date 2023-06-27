package utils

import (
	"log"
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

func TestFileGetContents(t *testing.T) {
	err := FilePutContents(contents, dst)
	assert.Equal(t, nil, err)
	_, err2 := FileGetContents(dst)
	assert.Equal(t, nil, err2)
	clear(err)
}

func TestFilePutContents(t *testing.T) {
	err := FilePutContents(contents, dst)
	assert.Equal(t, nil, err)
	clear(err)
}

func TestFileAppendContents(t *testing.T) {
	err := FilePutContents(contents, dst)
	assert.Equal(t, nil, err)
	err2 := FileAppendContents([]byte(`new line`), dst)
	assert.Equal(t, nil, err2)
	clear(err)
}

func TestFileExists(t *testing.T) {
	err := FilePutContents(contents, dst)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, FileExists(dst))
	clear(err)
}

func TestFileExt(t *testing.T) {
	err := FilePutContents(contents, dst)
	assert.Equal(t, nil, err)
	assert.Equal(t, "txt", FileExt(dst))
	clear(err)
}

func TestFileGetContentsReadLine(t *testing.T) {
	err := FilePutContents(contents, dst)
	assert.Equal(t, nil, err)

	ch, err2 := FileGetContentsReadLine(dst)
	assert.Equal(t, nil, err2)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for content := range ch {
			log.Printf("[TestFileGetContentsReadLine] content: %s\n", content)
			assert.NotEqual(t, "", content)
		}
	}()
	wg.Wait()
	clear(err)
}

func clear(err error) {
	if err == nil {
		_ = os.RemoveAll(dst)
	}
}
