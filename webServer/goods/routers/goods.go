package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"webServer/goods/api"
	"webServer/goods/middlewares"
)

func Routers(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	routerGroup.Use(middlewares.Trace())
	goodsGroup := routerGroup.Group("/goods")
	{
		g := api.NewGoods()
		goodsGroup.POST("", middlewares.JWTAuth(), middlewares.Authorities(), g.Create)
		goodsGroup.POST("/upload", middlewares.JWTAuth(), middlewares.Authorities(), g.Upload)
		goodsGroup.PATCH("/:id", middlewares.JWTAuth(), middlewares.Authorities(), g.UpdateStatus)
		goodsGroup.PUT("/:id", middlewares.JWTAuth(), middlewares.Authorities(), g.UpdateGoods)
		goodsGroup.DELETE("/:id", middlewares.JWTAuth(), middlewares.Authorities(), g.DeleteGoods)
		goodsGroup.GET("/img", g.PreviewImage)
		goodsGroup.GET("/:id", g.GetGoods)
		goodsGroup.GET("/:id/stocks", g.Stocks)
	}
	imgGroup := routerGroup.Group("upload")
	{
		imgGroup.GET("")
		imgGroup.POST("")
		imgGroup.StaticFS("/stativc", http.Dir("goods/"))
	}
	bannerGroup := routerGroup.Group("/banner")
	{
		b := api.NewBanner()
		bannerGroup.POST("", b.Add)
		bannerGroup.GET("", b.List)
		bannerGroup.PUT("/:id", b.Edit)
		bannerGroup.DELETE("/:id", b.Delete)
	}
	categoryGroup := routerGroup.Group("/category")
	{
		ct := api.NewCategory()
		categoryGroup.POST("", middlewares.JWTAuth(), middlewares.Authorities(), ct.Add)
		categoryGroup.GET("", ct.List)
		categoryGroup.PUT("/:id", middlewares.JWTAuth(), middlewares.Authorities(), ct.Edit)
		categoryGroup.DELETE("/:id", middlewares.JWTAuth(), middlewares.Authorities(), ct.Delete)
	}
	brandGroup := routerGroup.Group("/brand")
	{
		b := api.NewBrand()
		brandGroup.POST("", b.Add)
		brandGroup.GET("", b.List)
		brandGroup.PUT("/:id", b.Edit)
		brandGroup.DELETE("/:id", b.Delete)
	}
	categoryBrand := routerGroup.Group("/categorybrands")
	{
		br := api.NewBrand()
		categoryBrand.GET("/:id", br.GetBrandsByCategory)
		categoryBrand.GET("", br.CategoryBrandList)
		categoryBrand.POST("", br.CreateCategoryBrand)
		categoryBrand.PUT("/:id", br.UpdateCategoryBrand)
		categoryBrand.DELETE("/:id", br.DeleteCategoryBrand)
	}

}
