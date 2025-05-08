package utils

import (
	"bufio"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

// SaveFile2Dst 将文件保存到目标地址(拷贝文件)
//
// Deprecated: SaveFile2Dst is deprecated,it will be removed in the future.
//
// Please use File().Save2Dst() instead.
//
// file *multipart.FileHeader 文件
//
// dst string 拷贝到的目标地址
func SaveFile2Dst(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, src)

	return err
}

// FileGetContents 获取文件内容
//
// Deprecated: FileGetContents is deprecated,it will be removed in the future.
//
// Please use File().GetContents() instead.
//
// filename string 文件地址
func FileGetContents(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// FilePutContents 将内容写入文件(覆盖写)
//
// Deprecated: FilePutContents is deprecated,it will be removed in the future.
//
// Please use File().PutContents() instead.
//
// content []byte 写入的内容
//
// dst string 写入的目标地址
func FilePutContents(content []byte, dst string) error {
	return os.WriteFile(dst, content, 0644)
}

// FileAppendContents 将内容写入文件(追加写)
//
// Deprecated: FileAppendContents is deprecated,it will be removed in the future.
//
// Please use File().AppendContents() instead.
//
// content []byte 写入的内容
//
// dst string 写入的目标地址
func FileAppendContents(content []byte, dst string) error {
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	//defer func() {
	//	_ = f.Close()
	//}()
	_, err = f.Write(content)

	return err
}

// FileExists 检查文件上是否存在
//
// Deprecated: FileExists is deprecated,it will be removed in the future.
//
// Please use File().Exists() instead.
//
// dst string 目标地址
func FileExists(dst string) bool {
	ok, _ := FileExistsWithError(dst)

	return ok
}

// FileExistsWithError 检查文件上是否存在(会返回错误信息)
//
// Deprecated: FileExistsWithError is deprecated,it will be removed in the future.
//
// Please use File().ExistsWithError() instead.
//
// dst string 目标地址
func FileExistsWithError(dst string) (bool, error) {
	_, err := os.Stat(dst)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// FileExt 获取文件扩展名
//
// Deprecated: FileExt is deprecated,it will be removed in the future.
//
// Please use File().Ext() instead.
//
// 根据文件名最后一个.分隔来切分获取
func FileExt(filename string) string {
	filenameSplit := strings.Split(filename, ".")
	return filenameSplit[len(filenameSplit)-1]
}

// FileGetContentsReadLine 逐行读取文件内容
//
// Deprecated: FileGetContentsReadLine is deprecated,it will be removed in the future.
//
// Please use File().GetContentsReadLine() instead.
func FileGetContentsReadLine(dst string) (<-chan string, error) {
	ch := make(chan string)
	f, err := os.Open(dst)
	if err != nil {
		close(ch)
		return ch, err
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	go func() {
		for fileScanner.Scan() {
			ch <- fileScanner.Text()
		}
		close(ch)
	}()

	return ch, nil
}
