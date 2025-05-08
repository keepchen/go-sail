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
	// 创建一个缓冲区存储 multipart 数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加一个文件字段（例如字段名为 "file"）
	part, err := writer.CreateFormFile("file", "example.txt")
	if err != nil {
		panic(err)
	}

	// 写入文件内容
	content := []byte("Hello, World!")
	_, _ = part.Write(content)

	// 添加其他表单字段（可选）
	_ = writer.WriteField("key", "value")

	// 关闭 writer 完成数据写入
	_ = writer.Close()

	// 创建一个 HTTP 请求并解析 multipart 数据
	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 解析请求中的 multipart 数据
	_ = req.ParseMultipartForm(32 << 20) // 32MB 内存限制

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
