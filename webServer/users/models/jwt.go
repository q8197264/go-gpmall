package models

import (
	"github.com/golang-jwt/jwt"
)

//AuthorityId role
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint

	jwt.StandardClaims
}
