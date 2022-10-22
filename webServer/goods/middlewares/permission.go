package middlewares

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"webServer/goods/models"

	"github.com/gin-gonic/gin"
)

/**
获取 jwt 存储的数据
*/

// 普通用户授权
var generalAllowPaths = map[string]string{
	"/v1/goods": "get",
}

// 管理用户授权
var adminAllowPaths = map[string]string{
	"/v1/goods": "get|post|put|patch|delete",
}

func Authorities() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		if info, ok := claims.(*models.CustomClaims); ok {
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
		} else {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"code":   1,
				"errmsg": "jwt数据有误",
			})
			c.Abort()
			return
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
