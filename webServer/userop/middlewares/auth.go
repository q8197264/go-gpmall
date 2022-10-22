package middlewares

import "github.com/gin-gonic/gin"

// user role
func Authority() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Get("claims")
	}
}
