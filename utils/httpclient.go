package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	requestTimeout = time.Second * 5
)

type httpClientImpl struct {
}

type IHttpClient interface {
	// SendRequest 发送请求
	//
	// 默认Content-Type为application/json; charset=utf-8
	SendRequest(method, url string, jsonPayload []byte, headers map[string]string, timeout ...time.Duration) (response []byte, statusCode int, err error)
}

var _ IHttpClient = &httpClientImpl{}

// HttpClient 实例化http client工具类
func HttpClient() IHttpClient {
	return &httpClientImpl{}
}

// SendRequest 发送请求
//
// 默认Content-Type为application/json; charset=utf-8
func (httpClientImpl) SendRequest(method, url string, jsonPayload []byte, headers map[string]string, timeout ...time.Duration) (response []byte, statusCode int, err error) {
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
	response, err = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code: %d", resp.StatusCode)
	}

	fmt.Println("utils.httpclient.SendRequest", zap.String("body", string(response)))
	return
}
