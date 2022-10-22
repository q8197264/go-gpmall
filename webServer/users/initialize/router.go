package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"webServer/users/middlewares"
	"webServer/users/routers"
)

func Routes() *gin.Engine {
	router := gin.Default()
	RouterGroup := router.Group("/v1")
	RouterGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	RouterGroup.Use(middlewares.Cors())
	// router.LoadHTMLGlob("templates/*")
	RouterGroup.Use(middlewares.Trace())

	routers.InitUserRouter(RouterGroup)

	return router
}
