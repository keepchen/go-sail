package utils

import "golang.org/x/crypto/bcrypt"

type bcryptImpl struct{}

// IBCrypt bcrypt interface
type IBCrypt interface {
	// HashMake generates a hash string
	//
	// plaintext - plaintext string
	//
	// costs - cost factor, range ∈ [4,31], default: 10
	HashMake(plaintext string, costs ...int) (string, error)
	// HashCheck verifies a hash string
	//
	// plaintext - plaintext string
	//
	// hashed - hashed string
	HashCheck(plaintext, hashed string) (bool, error)
}

var _ IBCrypt = &bcryptImpl{}

// Bcrypt instantiates bcrypt hashing
func Bcrypt() IBCrypt {
	return &bcryptImpl{}
}

// HashMake generates a hash string
//
// plaintext - plaintext string
//
// costs - cost factor, range ∈ [4,31], default: 10
func (bcryptImpl) HashMake(plaintext string, costs ...int) (string, error) {
	var cost = bcrypt.DefaultCost
	if len(costs) > 0 {
		cost = costs[0]
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), cost)

	return string(bytes), err
}

// HashCheck verifies a hash string
//
// plaintext - plaintext string
//
// hashed - hashed string
func (bcryptImpl) HashCheck(plaintext, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))

	return err == nil, err
}
