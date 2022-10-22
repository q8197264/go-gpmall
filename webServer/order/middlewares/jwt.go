package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"webServer/order/global"
	"webServer/order/models"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"code":   405,
				"errmsg": "请登陆",
			})
			c.Abort()
			return
		}

		claims, err := parseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":   401,
				"errmsg": err.Error(),
			})
			c.Abort()
			return
		}

		// 授权成功, 保存用户信息
		c.Set("claims", claims)

		c.Next()
	}
}

func GenerateToken(claims *models.MyCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.ServerConfig.Jwt.Key))
}

func parseToken(tokenString string) (*models.MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(global.ServerConfig.Jwt.Key), nil
		},
	)
	if err != nil {
		if ve, b := err.(*jwt.ValidationError); b {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("That's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("Token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token is not active yet")
			} else {
				return nil, errors.New("Couldn't handle this token")
			}
		} else {
			return nil, err
		}
	}

	if claims, ok := token.Claims.(*models.MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("")
}

func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	claims, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}

	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
	return GenerateToken(claims)

}
