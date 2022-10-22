package initialize

import (
	"github.com/gin-gonic/gin"

	"webServer/inventory/routers"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	routers.Router(router)

	return router
}
