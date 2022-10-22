package routers

import (
	"webServer/userop/api"
	"webServer/userop/middlewares"

	"github.com/gin-gonic/gin"
)

func PostRouter(router *gin.RouterGroup) {
	postRouter := router.Group("post")
	postRouter.Use(middlewares.NewJwt())
	{
		p := api.NewPost()
		postRouter.GET("", p.List)
		postRouter.POST("", p.Add)
	}
}
