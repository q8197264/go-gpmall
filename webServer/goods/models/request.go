package models

import (
	"github.com/golang-jwt/jwt"
)

// jwt 鉴权数据
//AuthorityId role
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint

	jwt.StandardClaims
}
