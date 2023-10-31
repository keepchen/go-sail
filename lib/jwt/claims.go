package jwt

import (
	jwtLib "github.com/golang-jwt/jwt"
)

// AppClaims App票据声明
type AppClaims struct {
	Scopes   []string
	ScopeIDs []int64
	Name     string
	Email    string
	//more...
	jwtLib.StandardClaims
}

type MapClaims map[string]interface{}

func (c *MapClaims) Valid() error {
	//TODO valid fields
	return nil
}
