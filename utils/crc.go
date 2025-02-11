package utils

import (
	"hash/crc32"
	"hash/crc64"
)

type crc64Impl struct {
}

type ICrc64 interface {
	// Checksum 求crc64校验码
	Checksum(data []byte, table *crc64.Table) uint64
	// ChecksumECMA 求crc64校验码
	//
	// 使用ECMA多项式
	ChecksumECMA(data []byte) uint64
}

var _ ICrc64 = crc64Impl{}

// Crc64 实例化crc64工具类
func Crc64() ICrc64 {
	return &crc64Impl{}
}

// Checksum 求crc64校验码
func (crc64Impl) Checksum(data []byte, table *crc64.Table) uint64 {
	return crc64.Checksum(data, table)
}

// ChecksumECMA 求crc64校验码
//
// 使用ECMA多项式
func (crc64Impl) ChecksumECMA(data []byte) uint64 {
	crc64Table := crc64.MakeTable(crc64.ECMA)

	return crc64.Checksum(data, crc64Table)
}

type crc32Impl struct {
}

type ICrc32 interface {
	// Checksum 求crc32校验码
	Checksum(data []byte, table *crc32.Table) uint32
	// ChecksumIEEE 求crc32校验码
	//
	// 使用IEEE多项式
	ChecksumIEEE(data []byte) uint32
}

var _ ICrc32 = crc32Impl{}

// Crc32 实例化crc32工具类
func Crc32() ICrc32 {
	return &crc32Impl{}
}

// Checksum 求crc32校验码
func (crc32Impl) Checksum(data []byte, table *crc32.Table) uint32 {
	return crc32.Checksum(data, table)
}

// ChecksumIEEE 求crc32校验码
//
// 使用IEEE多项式
func (crc32Impl) ChecksumIEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
