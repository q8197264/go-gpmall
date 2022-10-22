package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"webServer/goods/forms"
	"webServer/goods/global"
	"webServer/goods/proto"
)

type category struct{}

func NewCategory() *category {
	c := &category{}
	return c
}

//添加类别
func (ct *category) Add(c *gin.Context) {
	var categoryForm forms.CategoryForm
	if err := c.ShouldBind(&categoryForm); err != nil {
		HandleValidateErrToHttp(err, c)
		return
	}

	// grpc 连接 goods-srv
	rsp, err := global.GoodsClient.CreateCategory(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryRequest{
			Name:     categoryForm.Name,
			ParentId: categoryForm.ParentId,
			Level:    categoryForm.Level,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    rsp,
		"message": "success",
	})
}

// 获取分类列表
func (ct *category) List(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	cid, _ := strconv.Atoi(id)

	rsp, err := global.GoodsClient.CategoryList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryFilterRequest{
			Id: int32(cid),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	r := make([]interface{}, 1)
	if err = json.Unmarshal([]byte(rsp.JsonData), &r); err != nil {
		zap.S().Errorw("[List]查询分类列表失败")
	}

	//分类数据 结果重构
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"data":     rsp.Data,
		"category": r,
		"message":  "success",
	})
}

//删除分类
func (ct *category) Delete(c *gin.Context) {
	// 连同所有子孙类一起删
	id := c.Param("id")
	cid, _ := strconv.Atoi(id)
	_, err := global.GoodsClient.DeleteCategory(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryByIdRequest{
			Id: int32(cid),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// 更新分类
func (ct *category) Edit(c *gin.Context) {
	var categoryForm forms.CategoryForm
	if errs := c.ShouldBind(&categoryForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}
	id := c.Param("id")
	cid, _ := strconv.Atoi(id)

	rsp, err := global.GoodsClient.UpdateCategory(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryRequest{
			Id:       int32(cid),
			Name:     categoryForm.Name,
			ParentId: categoryForm.ParentId,
			IsTab:    *categoryForm.IsTab,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    rsp,
		"message": "success",
	})
}
