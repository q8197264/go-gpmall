package routers

import (
	"webServer/userop/api"
	"webServer/userop/middlewares"

	"github.com/gin-gonic/gin"
)

func FavRouter(router *gin.RouterGroup) {
	favRouter := router.Group("/fav")
	favRouter.Use(middlewares.NewJwt())
	{
		fav := api.NewFav()
		favRouter.GET("", fav.List)
		favRouter.GET("detail/:gid", fav.Detail)
		favRouter.POST("", fav.Add)
		favRouter.DELETE("", fav.Delete)
	}
}
