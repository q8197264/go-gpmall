package middlewares

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"webServer/users/models"

	"github.com/gin-gonic/gin"
)

var generalAllowPaths = map[string]string{
	"/v1/user/login": "post",
}

var adminAllowPaths = map[string]string{
	"/v1/user/login": "post",
	"/v1/user":       "get",
}

func Authorities() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		info, _ := claims.(*models.CustomClaims)
		if info.AuthorityId > 0 {
			switch info.AuthorityId {
			case 1:
				if checkAccessPath(c, generalAllowPaths) {
					c.JSON(http.StatusOK, gin.H{
						"code":    0,
						"role":    1,
						"message": "普通用户",
					})
					c.Abort()
					return
				}
			case 2:
				if checkAccessPath(c, adminAllowPaths) {
					c.JSON(http.StatusOK, gin.H{
						"code":    0,
						"role":    2,
						"message": "管理员",
					})
					c.Abort()
					return
				}
			default:
			}
		}

		c.Next()
	}
}

func checkAccessPath(c *gin.Context, libs map[string]string) bool {
	forbid := true
	uri := c.Request.URL.String()
	if b := strings.Index(uri, "?"); b > 0 {
		uri = uri[:strings.Index(uri, "?")]
	}

	for k, v := range libs {
		ok, _ := regexp.MatchString(fmt.Sprintf(`(%s)`, v), strings.ToLower(c.Request.Method))
		if (strings.Index(uri, k) == 0) && ok {
			forbid = false
			break
		}
	}

	return forbid
}
