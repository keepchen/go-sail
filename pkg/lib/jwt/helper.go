package jwt

import (
	"errors"
	"os"
	"strings"
)

func fileGetContents(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func fileExists(dst string) bool {
	ok, _ := fileExistsWithError(dst)

	return ok
}

func fileExistsWithError(dst string) (bool, error) {
	_, err := os.Stat(dst)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func wordwrap(rawStr string, length int, split string) string {
	if len(rawStr) <= length || length < 1 {
		return rawStr
	}

	strSplit := strings.Split(rawStr, "")

	var (
		start    int
		end      = length
		finalStr string
	)

	for {
		if start > len(strSplit) {
			break
		}
		if end >= len(strSplit) {
			finalStr += strings.Join(strSplit[start:], "")
		} else {
			finalStr += strings.Join(strSplit[start:end], "") + split
		}

		start = end
		end += length
	}

	return finalStr
}
