package routers

import (
	"github.com/gin-gonic/gin"

	"webServer/order/api"
	"webServer/order/middlewares"
)

func Order(router *gin.RouterGroup) {
	orderRouter := router.Group("order").Use(middlewares.JwtAuth())
	{
		order := api.NewOrder()
		orderRouter.GET("", order.List)
		orderRouter.GET("/detail/:id", order.Detail)
		orderRouter.POST("", order.Create)
		orderRouter.DELETE("", order.Create)
	}

}
