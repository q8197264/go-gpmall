package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"webServer/userop/forms"
	"webServer/userop/global"
	"webServer/userop/proto"
)

type fav struct{}

func NewFav() *fav {
	return &fav{}
}

func (f *fav) List(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	rsp, err := global.GrpcFavClient.QueryFavList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserFavRequest{
			UserId: int32(claims.ID),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": rsp,
	})
}

func (f *fav) Detail(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	var favoritesForm forms.FavoritesForm
	if errs := c.ShouldBindUri(&favoritesForm); errs != nil {
		printValidateErrorTips(c, errs)
		return
	}
	rsp, err := global.GrpcFavClient.QueryFav(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserFavRequest{
			UserId:  int32(claims.ID),
			GoodsId: favoritesForm.GoodsId,
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"errmsg": rsp,
	})
}

func (f *fav) Add(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	var favoritesForm forms.FavoritesForm
	if err := c.ShouldBindJSON(&favoritesForm); err != nil {
		printValidateErrorTips(c, err)
		return
	}
	_, err = global.GrpcFavClient.AddFav(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserFavRequest{
			UserId:  int32(claims.ID),
			GoodsId: favoritesForm.GoodsId,
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"errmsg": "ok",
	})
}

func (f *fav) Delete(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	var favoritesForm forms.FavoritesForm
	if err := c.ShouldBindJSON(&favoritesForm); err != nil {
		printValidateErrorTips(c, err)
		return
	}
	_, err = global.GrpcFavClient.DeleteFav(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserFavRequest{
			UserId:  int32(claims.ID),
			GoodsId: favoritesForm.GoodsId,
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"errmsg": "ok",
	})
}
