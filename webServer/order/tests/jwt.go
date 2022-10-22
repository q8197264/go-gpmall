package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var mySigningKey = []byte("AllYourBase")

func main() {
	JwtAuth()
}

func JwtAuth() {
	// token := c.Request.Header.Get("o-token")
	token := GeneralToken()
	b := parserJwt(token)
	if b {
		// true
		// update
	} else {
		// false
		// 无效
		// login
	}
}

// Create the Claims
type MyCustomClaims struct {
	UID int    `json:"uid"`
	Foo string `json:"foo"`
	jwt.StandardClaims
}

func GeneralToken() string {
	var claims = MyCustomClaims{
		1,
		"sai",
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err.Error())
	}

	return ss
}

func updateJwt() {}

func parserJwt(tokenString string) bool {

	// sample token is expired.  override time so it parses as valid
	at(time.Unix(0, 0), func() {
		token, err := jwt.ParseWithClaims(
			tokenString,
			&MyCustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return mySigningKey, nil
			},
		)
		if token.Valid {
			fmt.Println("success")
			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				fmt.Printf("%v %v", claims.UID, claims.StandardClaims.ExpiresAt)
			} else {
				fmt.Println(err)
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("That's not even a token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				fmt.Println("Timing is everything")
			} else {
				fmt.Println("Couldn't handle this token:", err)
			}
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	})

	return true
}

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}
