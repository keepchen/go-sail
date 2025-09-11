package utils

import (
	"bufio"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type fileImpl struct {
}

type IFile interface {
	// Save2Dst 将文件保存到目标地址(拷贝文件)
	//
	// file *multipart.FileHeader 文件
	//
	// dst string 拷贝到的目标地址
	Save2Dst(file *multipart.FileHeader, dst string) error
	// GetContents 获取文件内容
	//
	// filename string 文件地址
	GetContents(filename string) ([]byte, error)
	//PutContents 将内容写入文件(覆盖写)
	//
	// content []byte 写入的内容
	//
	// dst string 写入的目标地址
	PutContents(content []byte, dst string) error
	// AppendContents 将内容写入文件(追加写)
	//
	// content []byte 写入的内容
	//
	// dst string 写入的目标地址
	AppendContents(content []byte, dst string) error
	// Exists 检查文件上是否存在
	//
	// dst string 目标地址
	Exists(dst string) bool
	// ExistsWithError 检查文件上是否存在(会返回错误信息)
	//
	// dst string 目标地址
	ExistsWithError(dst string) (bool, error)
	// Ext 获取文件扩展名
	//
	// 根据文件名最后一个.分隔来切分获取
	Ext(filename string) string
	// GetContentsReadLine 逐行读取文件内容
	GetContentsReadLine(dst string) (<-chan string, error)
}

var fi IFile = &fileImpl{}

// File 实例化file工具类
func File() IFile {
	return fi
}

// Save2Dst 将文件保存到目标地址(拷贝文件)
//
// file *multipart.FileHeader 文件
//
// dst string 拷贝到的目标地址
func (fileImpl) Save2Dst(file *multipart.FileHeader, dst string) error {
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

// GetContents 获取文件内容
//
// filename string 文件地址
func (fileImpl) GetContents(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// PutContents 将内容写入文件(覆盖写)
//
// content []byte 写入的内容
//
// dst string 写入的目标地址
func (fileImpl) PutContents(content []byte, dst string) error {
	return os.WriteFile(dst, content, 0644)
}

// AppendContents 将内容写入文件(追加写)
//
// content []byte 写入的内容
//
// dst string 写入的目标地址
func (fileImpl) AppendContents(content []byte, dst string) error {
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

// Exists 检查文件上是否存在
//
// dst string 目标地址
func (fileImpl) Exists(dst string) bool {
	ok, _ := File().ExistsWithError(dst)

	return ok
}

// ExistsWithError 检查文件上是否存在(会返回错误信息)
//
// dst string 目标地址
func (fileImpl) ExistsWithError(dst string) (bool, error) {
	_, err := os.Stat(dst)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// Ext 获取文件扩展名
//
// 根据文件名最后一个.分隔来切分获取
func (fileImpl) Ext(filename string) string {
	filenameSplit := strings.Split(filename, ".")
	return filenameSplit[len(filenameSplit)-1]
}

// GetContentsReadLine 逐行读取文件内容
func (fileImpl) GetContentsReadLine(dst string) (<-chan string, error) {
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
