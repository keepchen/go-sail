package utils

import "testing"

func TestHttpClientImplSendRequest(t *testing.T) {
	t.Log(HttpClient().SendRequest("GET", "https://github.com", nil, nil))
}
