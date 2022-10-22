package models

import "github.com/golang-jwt/jwt"

type MyCustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint

	jwt.StandardClaims
}
