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
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	err2 := File().AppendContents([]byte(`new line`), dst)
	assert.Equal(t, nil, err2)
	cleanClear(err)
}

func TestFileImplExists(t *testing.T) {
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, File().Exists(dst))
	cleanClear(err)
}

func TestFileImplExt(t *testing.T) {
	err := File().PutContents(contents, dst)
	assert.Equal(t, nil, err)
	assert.Equal(t, "txt", File().Ext(dst))
	cleanClear(err)
}

func TestFileImplGetContentsReadLine(t *testing.T) {
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
}

func cleanClear(err error) {
	if err == nil {
		_ = os.RemoveAll(dst)
	}
}
