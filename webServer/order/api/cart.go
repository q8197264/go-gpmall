package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"webServer/order/forms"
	"webServer/order/global"
	"webServer/order/proto"
)

type shopCart struct{}

func NewShopCart() *shopCart {
	return &shopCart{}
}

func (s *shopCart) Show(c *gin.Context) {
	claims, err := getJwtInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	rsp, err := global.ShopCartClient.QueryShopCart(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserInfoRequest{
			Id: int32(claims.ID),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rsp,
		"message": "success",
	})
}

func (s *shopCart) Add(c *gin.Context) {
	var addGoodsForm forms.AddGoodsForm
	if errs := c.ShouldBindJSON(&addGoodsForm); errs != nil {
		printValidateError(c, errs)
		return
	}

	claims, err := getJwtInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}
	_, err = global.ShopCartClient.AddGoodsToShopCart(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.ShopCartRequest{
			UserId:  int32(claims.ID),
			GoodsId: int32(addGoodsForm.GoodsId),
			Nums:    int32(addGoodsForm.Nums),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    claims,
		"message": "success",
	})
}

func (s *shopCart) Update(c *gin.Context) {
	var shopCartForm forms.AddGoodsForm
	if errs := c.ShouldBindJSON(&shopCartForm); errs != nil {
		printValidateError(c, errs)
		return
	}

	claims, err := getJwtInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	_, err = global.ShopCartClient.UpdateShopCart(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.ShopCartRequest{
			UserId:  int32(claims.ID),
			GoodsId: int32(shopCartForm.GoodsId),
			Nums:    int32(shopCartForm.Nums),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": claims,
	})
}

func (s *shopCart) Delete(c *gin.Context) {
	var id interface{}
	id = c.Param("id")
	gid, _ := id.(int32)

	claims, err := getJwtInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}
	_, err = global.ShopCartClient.DelGoodsInShopCart(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.ShopCartRequest{
			UserId:  int32(claims.ID),
			GoodsId: gid,
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    claims,
		"message": "success",
	})
}
