package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPImplGetLocal(t *testing.T) {
	ip, err := IP().GetLocal()
	t.Log("ip:", ip, "error:", err)
	assert.Equal(t, nil, err)
}
