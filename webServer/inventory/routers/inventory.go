package routers

import (
	"net/http"
	"webServer/inventory/api"
	"webServer/inventory/middlewares"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	router.Use(middlewares.Trace())
	routerGroup := router.Group("v1")
	inventoryGroup := routerGroup.Group("inventory")
	{
		inv := api.NewInventory()
		inventoryGroup.GET("/:id", inv.InvDetail)
		inventoryGroup.POST("", inv.SetInv)
	}
}
