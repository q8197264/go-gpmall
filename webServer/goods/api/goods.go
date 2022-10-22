package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"webServer/goods/forms"
	"webServer/goods/global"
	"webServer/goods/proto"
)

type goods struct{}

func NewGoods() *goods {
	g := &goods{}
	return g
}

// 创建商品
func (g *goods) Create(c *gin.Context) {
	var goodsForm forms.GoodsForm
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidateErrToHttp(err, c)
		return
	}
	rsp, err := global.GoodsClient.CreateGoods(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsRequest{
			CategoryId:  goodsForm.CategoryId,
			BrandId:     goodsForm.Brand,
			Name:        goodsForm.Name,
			GoodsSn:     goodsForm.GoodsSn,
			Subtitle:    goodsForm.Subtitle,
			MarketPrice: goodsForm.MarketPrice,
			ShopPrice:   goodsForm.ShopPrice,
			ShipFree:    *goodsForm.ShipFree,
			FrontImage:  goodsForm.FrontImage,
			Images:      goodsForm.Images,
			DescImages:  goodsForm.DescImages,
		},
	)
	if err != nil {
		HandleGrpcErrToHttp(err, c)
		return
	}

	// TODO: 库存服务

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    goodsForm,
		"message": rsp,
	})
}

// 上传文件
func (g *goods) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		zap.S().Infof("%s", err.Error())
		return
	}

	slide := []string{}
	desc := []string{}
	front := ""

	slidePath, descPath, frontPath := g.createImageDir()
	for fname, dir := range map[string]string{
		"slide_images[]": slidePath,
		"desc_images[]":  descPath,
		"front_image":    frontPath,
	} {
		files := form.File[fname]
		for _, file := range files {
			err := c.SaveUploadedFile(file, path.Join(dir, file.Filename))
			if err != nil {
				HandleSaveUploadFileErrToHttp(err, c)
				return
			} else {
				switch fname {
				case "slide_images[]":
					slide = append(slide, path.Join(dir, file.Filename))
				case "desc_images[]":
					desc = append(desc, path.Join(dir, file.Filename))
				case "front_image":
					front = path.Join(dir, file.Filename)
				default:
					// front = append(front, path.Join(dir, file.Filename))
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"front": front,
			"slide": slide,
			"desc":  desc,
		},
		"message": "success",
	})
}

// 创建图片目录
func (g *goods) createImageDir() (slide string, desc string, front string) {
	id := uuid.NewV4().String()
	path := fmt.Sprintf("images/goods/%d/%s", time.Now().Unix(), id)

	slide = filepath.Join(path, "slide")
	desc = filepath.Join(path, "desc")
	front = filepath.Join(path, "front")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		for _, dir := range []string{slide, desc, front} {
			os.MkdirAll(dir, os.ModePerm)
			os.Chmod(dir, 0777)
		}
	}

	return slide, desc, front
}

// 获取商品
func (g *goods) GetGoods(c *gin.Context) {
	var goodsByIdForm forms.GoodsByIdForm
	if err := c.ShouldBindUri(&goodsByIdForm); err != nil {
		HandleValidateErrToHttp(err, c)
		return
	}

	rsp, err := global.GoodsClient.GetGoodsDetail(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsByIdRequest{
			Id: goodsByIdForm.Id,
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

func (g *goods) PreviewImage(c *gin.Context) {
	src := c.DefaultQuery("src", "")
	file, err := ioutil.ReadFile(src)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查无此图片",
		})
		return
	}
	c.Writer.WriteString(string(file))
}

// 更新商品
func (g *goods) UpdateGoods(c *gin.Context) {
	var goodsForm forms.GoodsForm
	if errs := c.ShouldBind(&goodsForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}
	id := c.Param("id")
	gid, _ := strconv.ParseInt(id, 10, 32)
	_, err := global.GoodsClient.UpdateGoods(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsRequest{
			Id:          int32(gid),
			CategoryId:  goodsForm.CategoryId,
			BrandId:     goodsForm.Brand,
			Name:        goodsForm.Name,
			GoodsSn:     goodsForm.GoodsSn,
			Subtitle:    goodsForm.Subtitle,
			MarketPrice: goodsForm.MarketPrice,
			ShopPrice:   goodsForm.ShopPrice,
			ShipFree:    *goodsForm.ShipFree,
			FrontImage:  goodsForm.FrontImage,
			Images:      goodsForm.Images,
			DescImages:  goodsForm.DescImages,
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

func (g *goods) UpdateStatus(c *gin.Context) {
	var updateStatusForm forms.UpdateStatusForm
	if errs := c.ShouldBind(&updateStatusForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}

	id := c.Param("id")
	gid, _ := strconv.ParseInt(id, 10, 32)
	_, err := global.GoodsClient.UpdateStatus(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsRequest{
			Id:     int32(gid),
			OnSale: *updateStatusForm.On_Sale,
			IsNew:  *updateStatusForm.Is_New,
			IsHot:  *updateStatusForm.Is_Hot,
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

func (g *goods) Stocks(c *gin.Context) {
	// TODO: 库存
}

// 删除商品
func (g *goods) DeleteGoods(c *gin.Context) {
	var goodsByIdForm forms.GoodsByIdForm
	if errs := c.ShouldBindUri(&goodsByIdForm); errs != nil {
		HandleValidateErrToHttp(errs, c)
		return
	}
	_, err := global.GoodsClient.DeleteGoods(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsByIdRequest{
			Id: goodsByIdForm.Id,
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

// 参数验证错误提示
func HandleValidateErrToHttp(err error, c *gin.Context) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": ErrToJson(errs),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": "参数类型有误:" + err.Error(),
		})
	}
}

// 格式化错误信息
func ErrToJson(errs validator.ValidationErrors) map[string]string {
	tips := errs.Translate(global.Trans)
	data := make(map[string]string)
	var t []string
	for k, v := range tips {
		t = strings.Split(k, ".")
		data[t[len(t)-1]] = v
	}

	return data
}

// 保存上传图片错误
func HandleSaveUploadFileErrToHttp(err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"errmsg":  "upload fail",
			"message": err.Error(),
		})
	}
}

// grpc connect 错误
func HandleGrpcErrToHttp(err error, c *gin.Context) {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"code": 5, "errmsg": e.Message()})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "errmsg": e.Message()})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"code": 3, "errmsg": e.Message()})
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"code": 16, "errmsg": e.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "errmsg": e.Message()})
		}
		zap.S().Debugf("grpc错误: %s => %s", e.Code(), e.Message())
	}
}
