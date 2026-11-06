package utils

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

var (
	requestTimeout         = time.Second * 5
	defaultMaxReadBodySize = 1 << 21 //2MB
)

type httpClientImpl struct {
	maxReadBodySize int64
}

type IHttpClient interface {
	// WithLimit sets limit parameters
	//
	// maxReadBodySize maximum response body bytes to read
	WithLimit(maxReadBodySize int64) IHttpClient
	// SendRequest sends a request
	//
	// The default Content-Type is application/json; charset=utf-8
	//
	// # Note
	//
	// The response body read size is unlimited by default; use WithLimit to set a limit
	SendRequest(method, url string, jsonPayload []byte, headers map[string]string, timeout ...time.Duration) (response []byte, statusCode int, err error)
}

var hci IHttpClient = &httpClientImpl{}

// HttpClient instantiates HTTP client utilities
func HttpClient() IHttpClient {
	return hci
}

// WithLimit sets limit parameters
//
// maxReadBodySize maximum response body bytes to read
func (hc httpClientImpl) WithLimit(maxReadBodySize int64) IHttpClient {
	if maxReadBodySize <= 0 {
		maxReadBodySize = int64(defaultMaxReadBodySize)
	}
	hc.maxReadBodySize = maxReadBodySize

	return hc
}

// SendRequest sends a request
//
// The default Content-Type is application/json; charset=utf-8
//
// # Note
//
// The response body read size is unlimited by default; use WithLimit to set a limit
func (hc httpClientImpl) SendRequest(method, url string, jsonPayload []byte, headers map[string]string, timeout ...time.Duration) (response []byte, statusCode int, err error) {
	requestData := bytes.NewReader(jsonPayload)

	req, err := http.NewRequest(method, url, requestData)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	var reqTimeout time.Duration
	if len(timeout) == 0 || timeout[0] == 0 {
		reqTimeout = requestTimeout
	} else {
		reqTimeout = timeout[0]
	}

	client := &http.Client{Timeout: reqTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	statusCode = resp.StatusCode
	if hc.maxReadBodySize > 0 {
		response, err = io.ReadAll(io.LimitReader(resp.Body, hc.maxReadBodySize))
		return
	}
	response, err = io.ReadAll(resp.Body)
	//if resp.StatusCode != http.StatusOK {
	//	err = fmt.Errorf("http code: %d", resp.StatusCode)
	//}

	return
}
