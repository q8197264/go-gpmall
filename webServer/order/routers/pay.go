package routers

import (
	"webServer/order/api"

	"github.com/gin-gonic/gin"
)

func Pay(router *gin.RouterGroup) {
	alipayRouter := router.Group("/alipay")
	{
		a := api.NewAlipay()
		alipayRouter.GET("/callback", a.Callback)
		alipayRouter.POST("/notify", a.Notify)
	}
}
