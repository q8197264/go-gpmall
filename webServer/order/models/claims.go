package models

import "github.com/golang-jwt/jwt"

// Create the Claims
type MyCustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint

	jwt.StandardClaims
}
