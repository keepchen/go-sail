package utils

import "testing"

func TestSendRequest(t *testing.T) {
	t.Log(SendRequest("GET", "https://github.com", nil, nil))
}
