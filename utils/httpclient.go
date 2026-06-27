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
	// WithLimit 设置限制参数
	//
	// maxReadBodySize 最大读取响应体字节数
	WithLimit(maxReadBodySize int64) IHttpClient
	// SendRequest 发送请求
	//
	// 默认Content-Type为application/json; charset=utf-8
	//
	// # 注意
	//
	// 默认读取的响应体大小为2MB，若要改变此限制，请使用 WithLimit 进行设置
	SendRequest(method, url string, jsonPayload []byte, headers map[string]string, timeout ...time.Duration) (response []byte, statusCode int, err error)
}

var hci IHttpClient = &httpClientImpl{}

// HttpClient 实例化http client工具类
func HttpClient() IHttpClient {
	return hci
}

// WithLimit 设置限制参数
//
// maxReadBodySize 最大读取响应体字节数
func (hc httpClientImpl) WithLimit(maxReadBodySize int64) IHttpClient {
	if maxReadBodySize <= 0 {
		maxReadBodySize = int64(defaultMaxReadBodySize)
	}
	hc.maxReadBodySize = maxReadBodySize

	return hc
}

// SendRequest 发送请求
//
// 默认Content-Type为application/json; charset=utf-8
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
	response, err = io.ReadAll(io.LimitReader(resp.Body, hc.maxReadBodySize))
	//if resp.StatusCode != http.StatusOK {
	//	err = fmt.Errorf("http code: %d", resp.StatusCode)
	//}

	return
}
