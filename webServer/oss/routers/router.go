package routers

import (
	"net/http"
	"webServer/oss/api"
	"webServer/oss/middlewares"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	router.Use(middlewares.Trace())
	routerGroup := router.Group("v1")
	fileRouter := routerGroup.Group("qiniu")
	{
		qiniu := api.NewQiniu()
		fileRouter.GET("", qiniu.GetUpToken)
		fileRouter.POST("", qiniu.Notify)
	}
}
