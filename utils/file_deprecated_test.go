package utils

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	contentsDeprecated = []byte(`hello
world
hello
go-sail`)
	dstDeprecated = "./file_deprecated_test_will_be_deleted_after_testcase.txt"
)

func TestFileGetContents(t *testing.T) {
	err := FilePutContents(contentsDeprecated, dstDeprecated)
	assert.Equal(t, nil, err)
	_, err2 := FileGetContents(dstDeprecated)
	assert.Equal(t, nil, err2)
	clearDeprecated(err)
}

func TestFilePutContents(t *testing.T) {
	err := FilePutContents(contentsDeprecated, dstDeprecated)
	assert.Equal(t, nil, err)
	clearDeprecated(err)
}

func TestFileAppendContents(t *testing.T) {
	err := FilePutContents(contentsDeprecated, dstDeprecated)
	assert.Equal(t, nil, err)
	err2 := FileAppendContents([]byte(`new line`), dstDeprecated)
	assert.Equal(t, nil, err2)
	clearDeprecated(err)
}

func TestFileExists(t *testing.T) {
	err := FilePutContents(contentsDeprecated, dstDeprecated)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, FileExists(dstDeprecated))
	clearDeprecated(err)
}

func TestFileExt(t *testing.T) {
	err := FilePutContents(contentsDeprecated, dstDeprecated)
	assert.Equal(t, nil, err)
	assert.Equal(t, "txt", FileExt(dstDeprecated))
	clearDeprecated(err)
}

func TestFileGetContentsReadLine(t *testing.T) {
	err := FilePutContents(contentsDeprecated, dstDeprecated)
	assert.Equal(t, nil, err)

	ch, err2 := FileGetContentsReadLine(dstDeprecated)
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
	clearDeprecated(err)
}

func clearDeprecated(err error) {
	if err == nil {
		_ = os.RemoveAll(dstDeprecated)
	}
}
