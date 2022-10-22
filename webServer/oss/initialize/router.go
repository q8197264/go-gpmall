package initialize

import (
	"webServer/oss/middlewares"
	"webServer/oss/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Use(middlewares.Cors())
	routers.Router(router)

	return router
}
