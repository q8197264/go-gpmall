package middlewares

import "github.com/gin-gonic/gin"

func Authority() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Get("claims")
	}
}
