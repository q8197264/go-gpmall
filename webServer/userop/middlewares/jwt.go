package middlewares

import (
	"errors"
	"net/http"
	"time"
	"webServer/userop/global"
	"webServer/userop/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var signingKey = []byte(global.ServerConfig.Jwt.Key)

func NewJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("x-token")
		if tokenString == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":   1,
				"errmsg": "请登陆",
			})
			c.Abort()
			return
		}
		claims, err := paseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":   1,
				"errmsg": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func paseToken(tokenString string) (*models.MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(global.ServerConfig.Jwt.Key), nil
		},
	)
	if err != nil {
		if err, ok := err.(*jwt.ValidationError); ok {
			if err.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("That's not even a token")
			} else if err.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("Token is expired")
			} else if err.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token not active yet")
			}
			return nil, err
		}
	}
	claims, ok := token.Claims.(*models.MyCustomClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, err
}

func GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.ServerConfig.Jwt.Key))
}

func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	claims, err := paseToken(tokenString)
	if err != nil {
		return "", err
	}
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(24 * 30 * time.Hour).Unix()
	return GenerateToken(claims)
}
