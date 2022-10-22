package routers

import (
	"webServer/userop/api"
	"webServer/userop/middlewares"

	"github.com/gin-gonic/gin"
)

func AddressRouter(router *gin.RouterGroup) {
	addressRouter := router.Group("address")
	addressRouter.Use(middlewares.NewJwt())
	{
		a := api.NewAddress()
		addressRouter.GET("", a.List)
		addressRouter.POST("", a.Add)
		addressRouter.PUT("/:id", a.Update)
		addressRouter.DELETE("/:id", a.Delete)
	}
}
