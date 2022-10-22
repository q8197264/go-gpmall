package middlewares

import (
	"errors"
	"net/http"
	"time"
	"webServer/goods/global"
	"webServer/goods/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/*
Gin 的中间件
  获取t-token并解析jwt数据

*/
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"message": "请登陆",
			})
			c.Abort()
			return
		}

		//解析 token
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "授权已过期",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "未授权:" + err.Error(),
			})
			c.Abort()
			return
		}

		// 授权成功, 保存用户信息
		c.Set("claims", claims)

		c.Next()
	}
}

/**
	jwt driver
**/
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalFormed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTConfig.SigningKey), //可以设置过期时间
	}
}

// 生成 token
func (j *JWT) GenerateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenStr string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalFormed
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, v
			}
		}
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// 刷新 token
func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.GenerateToken(*claims)
	}

	return "", nil
}
