package utils

import (
	"hash/crc32"
	"hash/crc64"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	crcRawData    = []byte(`hello world!你好！こんにちは！안녕하세요!`)
	crc32checksum = uint32(411352257)
	crc64checksum = uint64(12645825114846618640)
)

func TestCrc32ImplChecksumIEEE(t *testing.T) {
	checksum := Crc32().ChecksumIEEE(crcRawData)
	t.Log(checksum)
	assert.Equal(t, crc32checksum, checksum)
}

func TestCrc32ImplChecksum(t *testing.T) {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := Crc32().Checksum(crcRawData, table)
	t.Log(checksum)
	assert.Equal(t, crc32checksum, checksum)
}

func TestCrc64ImplChecksumECMA(t *testing.T) {
	checksum := Crc64().ChecksumECMA(crcRawData)
	t.Log(checksum)
	assert.Equal(t, crc64checksum, checksum)
}

func TestCrc64ImplChecksum(t *testing.T) {
	table := crc64.MakeTable(crc64.ECMA)
	checksum := Crc64().Checksum(crcRawData, table)
	t.Log(checksum)
	assert.Equal(t, crc64checksum, checksum)
}
