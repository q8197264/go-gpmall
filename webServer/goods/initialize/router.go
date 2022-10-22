package initialize

import (
	"webServer/goods/middlewares"
	"webServer/goods/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8M
	router.Use(middlewares.Cors())

	routerGroup := router.Group("/v1")
	routers.Routers(routerGroup)

	return router
}
