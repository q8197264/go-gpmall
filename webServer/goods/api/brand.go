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

type brand struct{}

func NewBrand() *brand {
	b := &brand{}
	return b
}

// 注意与分类的关系
func (b *brand) List(c *gin.Context) {
	p := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(p)
	l := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(l)
	rsp, err := global.GoodsClient.GetBrandList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BrandFilterRquest{
			Page:  int32(page),
			Limit: int32(limit),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"total":   rsp.Total,
		"data":    rsp.Data,
		"message": "success",
	})
}

func (b *brand) Add(c *gin.Context) {
	var brandForm forms.BrandForm
	if errs := c.ShouldBind(&brandForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}

	rsp, err := global.GoodsClient.CreateBrand(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CreateBrandRequest{
			Name: brandForm.Name,
			Logo: brandForm.Logo,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	res := make(map[string]interface{})
	res["id"] = rsp.Id
	res["name"] = rsp.Name
	res["logo"] = rsp.Logo

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    res,
		"message": "success",
	})
}

func (b *brand) Edit(c *gin.Context) {
	var brandForm forms.BrandForm
	if errs := c.ShouldBind(&brandForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}
	bid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsClient.UpdateBrand(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BrandInfoRequest{
			Id:   int32(bid),
			Name: brandForm.Name,
			Logo: brandForm.Logo,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (b *brand) Delete(c *gin.Context) {
	bid, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsClient.DeleteBrand(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.BrandByIdRequest{
			Id: int32(bid),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 分类下创建品牌
func (b *brand) CreateCategoryBrand(c *gin.Context) {
	var categoryBrandForm forms.CategoryBrandForm
	if errs := c.ShouldBind(&categoryBrandForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}
	_, err := global.GoodsClient.CreateCategoryBrand(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryBrandRequest{
			CategoryId: categoryBrandForm.CategoryId,
			BrandId:    categoryBrandForm.BrandId,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 获取指定分类下所有品牌
func (b *brand) GetBrandsByCategory(c *gin.Context) {
	id := c.Param("id")
	cid, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	rsp, err := global.GoodsClient.GetBrandsByCategory(
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
		"code":    0,
		"data":    rsp,
		"message": "success",
	})
}

func (b *brand) CategoryBrandList(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	rsp, err := global.GoodsClient.CategoryBrandList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryBrandFilterRequest{
			Page:  int32(page),
			Limit: int32(limit),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}
	res := make([]interface{}, 0)
	for _, v := range rsp.Data {
		data := make(map[string]interface{})
		data["id"] = v.Id
		data["category"] = map[string]interface{}{
			"id":   v.Category.Id,
			"name": v.Category.Name,
		}
		data["brand"] = map[string]interface{}{
			"id":   v.Brand.Id,
			"name": v.Brand.Name,
			"logo": v.Brand.Logo,
		}
		res = append(res, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"total":   rsp.Total,
		"data":    res,
		"message": "success",
	})
}

func (b *brand) UpdateCategoryBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var categoryBrandForm forms.CategoryBrandForm
	if errs := c.ShouldBind(&categoryBrandForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}

	_, err = global.GoodsClient.UpdateCategoryBrand(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryBrandRequest{
			Id:         int32(id),
			CategoryId: categoryBrandForm.CategoryId,
			BrandId:    categoryBrandForm.BrandId,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    categoryBrandForm,
		"message": "success",
	})
}

func (b *brand) DeleteCategoryBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsClient.DeleteCategoryBrand(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CategoryBrandRequest{
			Id: int32(id),
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
