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
	// Save2Dst saves a file to the target path (copies the file)
	//
	// file *multipart.FileHeader file
	//
	// dst string target path for the copy
	Save2Dst(file *multipart.FileHeader, dst string) error
	// GetContents retrieves file content
	//
	// filename string file path
	GetContents(filename string) ([]byte, error)
	//PutContents writes content to a file (overwrites)
	//
	// content []byte content to write
	//
	// dst string target path to write to
	PutContents(content []byte, dst string) error
	// AppendContents writes content to a file (appends)
	//
	// content []byte content to write
	//
	// dst string target path to write to
	AppendContents(content []byte, dst string) error
	// Exists checks whether a file exists
	//
	// dst string target path
	Exists(dst string) bool
	// ExistsWithError checks whether a file exists (returns error information)
	//
	// dst string target path
	ExistsWithError(dst string) (bool, error)
	// Ext retrieves the file extension
	//
	// Split at the last dot in the filename
	Ext(filename string) string
	// GetContentsReadLine reads file content line by line
	GetContentsReadLine(dst string) (<-chan string, error)
}

var fi IFile = &fileImpl{}

// File instantiates file utilities
func File() IFile {
	return fi
}

// Save2Dst saves a file to the target path (copies the file)
//
// file *multipart.FileHeader file
//
// dst string target path for the copy
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

// GetContents retrieves file content
//
// filename string file path
func (fileImpl) GetContents(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// PutContents writes content to a file (overwrites)
//
// content []byte content to write
//
// dst string target path to write to
func (fileImpl) PutContents(content []byte, dst string) error {
	return os.WriteFile(dst, content, 0644)
}

// AppendContents writes content to a file (appends)
//
// content []byte content to write
//
// dst string target path to write to
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

// Exists checks whether a file exists
//
// dst string target path
func (fileImpl) Exists(dst string) bool {
	ok, _ := File().ExistsWithError(dst)

	return ok
}

// ExistsWithError checks whether a file exists (returns error information)
//
// dst string target path
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

// Ext retrieves the file extension
//
// Split at the last dot in the filename
func (fileImpl) Ext(filename string) string {
	filenameSplit := strings.Split(filename, ".")
	return filenameSplit[len(filenameSplit)-1]
}

// GetContentsReadLine reads file content line by line
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
