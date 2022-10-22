package initialize

import (
	"net/http"
	"webServer/userop/middlewares"
	"webServer/userop/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()
	router := engine.Group("v1")
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": true,
		})
	})
	router.Use(middlewares.Trace())
	routers.FavRouter(router)
	routers.AddressRouter(router)
	routers.PostRouter(router)

	return engine
}
