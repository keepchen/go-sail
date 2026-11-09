package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Config controls Argon2id hashing. Memory must be at least 8 KiB per
// lane and at most 128 MiB, with 1–10 iterations and 1–16 lanes. The product
// of Memory and Iterations must not exceed 256 MiB-passes. HashMake accepts
// 8–64 byte salts and 4–64 byte keys.
type Argon2Config struct {
	Memory      uint32 // Memory usage (KiB)
	Iterations  uint32 // Iteration count
	Parallelism uint8  // Parallelism
	SaltLength  uint32 // Salt length
	KeyLength   uint32 // Output key length
}

type argon2Impl struct {
	config Argon2Config
}

// IArgon2 Argon2 interface
type IArgon2 interface {
	// HashMake generates a hash string
	//
	// plaintext - plaintext string
	HashMake(plaintext string) (string, error)
	// HashCheck verifies a hash string
	//
	// plaintext - plaintext string
	//
	// hashed - hashed string
	HashCheck(plaintext, hashed string) (bool, error)
}

var _ IArgon2 = &argon2Impl{}

// Default Argon2 configuration
var defaultArgon2Config = Argon2Config{
	Memory:      64 * 1024, //64MB
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

const (
	argon2EncodeFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	argon2DecodeFormat = "m=%d,t=%d,p=%d"

	// Per-call limits bound the work requested by a stored hash.
	argon2MaxMemory        = 128 * 1024 // KiB
	argon2MaxIterations    = 10
	argon2MaxParallelism   = 16
	argon2MaxMemoryPasses  = 256 * 1024 // KiB times iterations
	argon2MinNewSaltLength = 8
	argon2MaxSaltLength    = 64
	argon2MinKeyLength     = 4
	argon2MaxKeyLength     = 64
	argon2MaxEncodedLength = 256 // Includes the header and both Base64 fields.
)

func (c Argon2Config) validate() error {
	if c.Parallelism == 0 || c.Parallelism > argon2MaxParallelism {
		return fmt.Errorf("Argon2 parallelism must be between 1 and %d", argon2MaxParallelism)
	}
	minMemory := 8 * uint32(c.Parallelism)
	if c.Memory < minMemory || c.Memory > argon2MaxMemory {
		return fmt.Errorf("Argon2 memory must be between %d and %d KiB", minMemory, argon2MaxMemory)
	}
	if c.Iterations == 0 || c.Iterations > argon2MaxIterations {
		return fmt.Errorf("Argon2 iterations must be between 1 and %d", argon2MaxIterations)
	}
	if c.Memory*c.Iterations > argon2MaxMemoryPasses {
		return fmt.Errorf("Argon2 memory times iterations must not exceed %d KiB-passes", argon2MaxMemoryPasses)
	}
	if c.SaltLength == 0 || c.SaltLength > argon2MaxSaltLength {
		return fmt.Errorf("Argon2 salt length must be between 1 and %d bytes", argon2MaxSaltLength)
	}
	if c.KeyLength < argon2MinKeyLength || c.KeyLength > argon2MaxKeyLength {
		return fmt.Errorf("Argon2 key length must be between %d and %d bytes", argon2MinKeyLength, argon2MaxKeyLength)
	}
	return nil
}

// Argon2 instantiates the Argon2 algorithm
func Argon2(config ...Argon2Config) IArgon2 {
	if len(config) > 0 {
		return &argon2Impl{config[0]}
	}
	return &argon2Impl{
		config: defaultArgon2Config,
	}
}

// HashMake generates a hash string
//
// plaintext - plaintext string
func (a *argon2Impl) HashMake(plaintext string) (string, error) {
	if a.config.SaltLength < argon2MinNewSaltLength {
		return "", fmt.Errorf("Argon2 salt length must be at least %d bytes", argon2MinNewSaltLength)
	}
	if err := a.config.validate(); err != nil {
		return "", err
	}

	salt := make([]byte, a.config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hashed := argon2.IDKey([]byte(plaintext), salt,
		a.config.Iterations,
		a.config.Memory,
		a.config.Parallelism,
		a.config.KeyLength)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hashed)

	result := fmt.Sprintf(argon2EncodeFormat, argon2.Version, a.config.Memory, a.config.Iterations, a.config.Parallelism, b64salt, b64hash)

	return result, nil
}

// HashCheck verifies an Argon2id hash and rejects malformed or out-of-range
// parameters before deriving the key. It accepts legacy non-empty salts
// shorter than the minimum used by HashMake.
//
// plaintext - plaintext string
//
// hashed - hashed string
func (a *argon2Impl) HashCheck(plaintext, hashed string) (bool, error) {
	if len(hashed) > argon2MaxEncodedLength {
		return false, fmt.Errorf("Argon2 hash exceeds %d bytes", argon2MaxEncodedLength)
	}
	parts := strings.Split(hashed, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" || strings.ContainsAny(hashed, "\r\n") {
		return false, fmt.Errorf("invalid Argon2id hash format")
	}
	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		return false, fmt.Errorf("unsupported Argon2 version")
	}

	var cfg Argon2Config
	_, err := fmt.Sscanf(parts[3], argon2DecodeFormat, &cfg.Memory, &cfg.Iterations, &cfg.Parallelism)
	if err != nil || parts[3] != fmt.Sprintf(argon2DecodeFormat, cfg.Memory, cfg.Iterations, cfg.Parallelism) {
		return false, fmt.Errorf("invalid Argon2 cost parameters")
	}

	encoding := base64.RawStdEncoding.Strict()
	salt, err := encoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := encoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	cfg.SaltLength = uint32(len(salt))
	cfg.KeyLength = uint32(len(hash))
	if err := cfg.validate(); err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(plaintext), salt,
		cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	// Use a constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
