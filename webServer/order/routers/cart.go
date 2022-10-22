package routers

import (
	"webServer/order/api"
	"webServer/order/middlewares"

	"github.com/gin-gonic/gin"
)

func Cart(router *gin.RouterGroup) {
	cartRouter := router.Group("cart").Use(middlewares.JwtAuth())
	{
		cart := api.NewShopCart()
		cartRouter.GET("", cart.Show)
		cartRouter.POST("", cart.Add)
		cartRouter.PUT("", cart.Update)
		cartRouter.DELETE("/:id", cart.Delete)
	}
}
