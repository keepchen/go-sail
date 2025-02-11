package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	t.Log("ip:", ip, "error:", err)
	assert.Equal(t, nil, err)
}
