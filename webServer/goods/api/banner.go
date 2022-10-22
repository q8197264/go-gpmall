package api

import (
	"context"
	"net/http"
	"strconv"
	"webServer/goods/forms"
	"webServer/goods/global"
	"webServer/goods/proto"

	"github.com/gin-gonic/gin"
)

type banner struct{}

func NewBanner() *banner {
	b := &banner{}
	return b
}

func (b *banner) Add(c *gin.Context) {
	var bannerForm forms.BannerForm
	if errs := c.ShouldBind(&bannerForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}

	rsp, err := global.GoodsClient.CreateBanner(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BannerRequest{
			Index: bannerForm.Index,
			Image: bannerForm.Image,
			Url:   bannerForm.Url,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rsp,
		"message": "success",
	})
}

func (b *banner) List(c *gin.Context) {
	rsp, err := global.GoodsClient.GetBannerList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BannerFilterRequest{
			Page:  1,
			Limit: 10,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rsp,
	})
}

func (b *banner) Edit(c *gin.Context) {
	var bannerForm forms.BannerForm
	if errs := c.ShouldBind(&bannerForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}
	id := c.Param("id")
	bid, _ := strconv.Atoi(id)
	rsp, err := global.GoodsClient.UpdateBanner(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BannerRequest{
			Id:    int32(bid),
			Index: bannerForm.Index,
			Image: bannerForm.Image,
			Url:   bannerForm.Url,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": rsp,
	})
}

func (b *banner) Delete(c *gin.Context) {
	id := c.Param("id")
	bid, _ := strconv.Atoi(id)
	rsp, err := global.GoodsClient.DeleteBanner(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BannerByIdRequest{
			Id: int32(bid),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": rsp,
	})
}
