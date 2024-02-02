package utils

import (
	"hash/crc32"
	"hash/crc64"
)

// Crc64Checksum 求crc64校验码
func Crc64Checksum(data []byte, table *crc64.Table) uint64 {
	return crc64.Checksum(data, table)
}

// Crc64ChecksumECMA 求crc64校验码
//
// 使用ECMA多项式
func Crc64ChecksumECMA(data []byte) uint64 {
	crc64Table := crc64.MakeTable(crc64.ECMA)

	return crc64.Checksum(data, crc64Table)
}

// Crc32Checksum 求crc32校验码
func Crc32Checksum(data []byte, table *crc32.Table) uint32 {
	return crc32.Checksum(data, table)
}

// Crc32ChecksumIEEE 求crc32校验码
//
// 使用IEEE多项式
func Crc32ChecksumIEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
