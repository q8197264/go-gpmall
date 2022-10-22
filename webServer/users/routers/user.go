package routers

import (
	"github.com/gin-gonic/gin"

	"webServer/users/api"
	"webServer/users/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	routerGroup := Router.Group("/user")
	{
		u := api.NewUser()

		routerGroup.GET("", middlewares.JWTAuth(), middlewares.Authorities(), u.GetUserList)
		routerGroup.PUT("", middlewares.JWTAuth(), u.UpdateUserInfo)
		routerGroup.POST("", u.CreateUser)
		routerGroup.POST("/login", u.Login)
		routerGroup.GET("/:uid", middlewares.JWTAuth(), middlewares.Authorities(), u.GetUserById)
		routerGroup.GET("/mobile", middlewares.JWTAuth(), middlewares.Authorities(), u.GetUserInfo)

		routerGroup.GET("/captcha", api.GetCaptcha)
		routerGroup.POST("/captcha", api.CaptchaVerify)
		// routerGroup.DELETE("", )
	}
}
