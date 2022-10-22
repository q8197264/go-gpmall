package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"webServer/order/middlewares"
	"webServer/order/routers"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()

	group := engine.Group("v1")
	group.Use(middlewares.Cors())
	group.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": true,
		})
	})
	group.Use(middlewares.Trace())
	routers.Order(group)
	routers.Cart(group)
	routers.Pay(group)

	return engine
}
