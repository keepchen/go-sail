package utils

import (
	"hash/crc32"
	"hash/crc64"
)

type crc64Impl struct {
}

type ICrc64 interface {
	// Checksum calculates a CRC64 checksum
	Checksum(data []byte, table *crc64.Table) uint64
	// ChecksumECMA calculates a CRC64 checksum
	//
	// Uses the ECMA polynomial
	ChecksumECMA(data []byte) uint64
}

var c64i ICrc64 = &crc64Impl{}

// Crc64 instantiates CRC64 utilities
func Crc64() ICrc64 {
	return c64i
}

// Checksum calculates a CRC64 checksum
func (crc64Impl) Checksum(data []byte, table *crc64.Table) uint64 {
	return crc64.Checksum(data, table)
}

// ChecksumECMA calculates a CRC64 checksum
//
// Uses the ECMA polynomial
func (crc64Impl) ChecksumECMA(data []byte) uint64 {
	crc64Table := crc64.MakeTable(crc64.ECMA)

	return crc64.Checksum(data, crc64Table)
}

type crc32Impl struct {
}

type ICrc32 interface {
	// Checksum calculates a CRC32 checksum
	Checksum(data []byte, table *crc32.Table) uint32
	// ChecksumIEEE calculates a CRC32 checksum
	//
	// Uses the IEEE polynomial
	ChecksumIEEE(data []byte) uint32
}

var c32i ICrc32 = crc32Impl{}

// Crc32 instantiates CRC32 utilities
func Crc32() ICrc32 {
	return c32i
}

// Checksum calculates a CRC32 checksum
func (crc32Impl) Checksum(data []byte, table *crc32.Table) uint32 {
	return crc32.Checksum(data, table)
}

// ChecksumIEEE calculates a CRC32 checksum
//
// Uses the IEEE polynomial
func (crc32Impl) ChecksumIEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
