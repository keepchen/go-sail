package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	t.Log("ip:", ip, "error:", err)
	assert.Equal(t, nil, err)
}
