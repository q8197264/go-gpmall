package api

import (
	"context"
	"net/http"
	"strings"
	"webServer/inventory/forms"
	"webServer/inventory/global"
	"webServer/inventory/proto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventory struct{}

func NewInventory() *inventory {
	i := &inventory{}
	return i
}

func (i *inventory) SetInv(c *gin.Context) {
	var inventoryForm forms.InventoryForm
	if errs := c.ShouldBindJSON(&inventoryForm); errs != nil {
		ValidateFormField(c, errs)
		return
	}
	_, err := global.GrpcClient.SetInv(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsInvInfo{
			GoodsId: int32(inventoryForm.GoodsId),
			Num:     int32(inventoryForm.Nums),
		},
	)
	if err != nil {
		PrintGrpcErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (i *inventory) InvDetail(c *gin.Context) {
	var invDetailForm forms.InvDetailForm
	if errs := c.ShouldBindUri(&invDetailForm); errs != nil {
		ValidateFormField(c, errs)
		return
	}
	rsp, err := global.GrpcClient.InvDetail(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.GoodsInvInfo{
			GoodsId: int32(invDetailForm.GoodsId),
		},
	)
	if err != nil {
		PrintGrpcErrors(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": rsp,
	})
}

// func (i *inventory) Sell(c *gin.Context) {
// 	var inventoryForm forms.InventoryForm
// 	if err := c.ShouldBindJSON(&inventoryForm); err != nil {
// 		ValidateFormField(c, err)
// 		return
// 	}

// 	rsp, err := global.GrpcClient.Sell(context.Background(), &proto.SellInfo{})
// }

// func (i *inventory) Reback(c *gin.Context) {

// }

func ValidateFormField(c *gin.Context, errs error) {
	if err, ok := errs.(validator.ValidationErrors); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": formatErrorAsJson(err),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":   2,
			"errmsg": "参数类型有误",
		})
	}
}

func formatErrorAsJson(err validator.ValidationErrors) map[string]string {
	msg := err.Translate(*global.Trans)
	data := make(map[string]string)
	for k, v := range msg {
		t := strings.Split(k, ".")
		data[t[len(t)-1]] = v
	}

	return data
}

func PrintGrpcErrors(c *gin.Context, err error) {
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
