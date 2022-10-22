package api

import (
	"context"
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"webServer/order/forms"
	"webServer/order/global"
	"webServer/order/proto"
	"webServer/order/response"
)

type order struct{}

func NewOrder() *order {
	return &order{}
}

func OrderResponse(rsp *proto.OrderDetailResponse) response.OrderDetail {
	order := response.OrderDetail{
		OrderId:      rsp.Id,
		OrderSn:      rsp.OrderSn,
		UserId:       rsp.UserId,
		SignerMobile: rsp.SignerMobile,
		Amount:       rsp.Amount,
		PayAmount:    rsp.PayAmount,
		PayTime:      rsp.PayTime,
	}
	for _, item := range rsp.Goods {
		r := response.OrderGoodsDetail{}
		r.GoodsId = item.GoodsId
		r.GoodsName = item.GoodsName
		r.MarketPrice = item.MarketPrice
		r.ShopPrice = item.ShopPrice
		r.Nums = item.Nums
		order.GoodsList = append(order.GoodsList, r)
	}

	return order
}

func (o *order) Create(c *gin.Context) {
	var orderInfoForm forms.OrderInfoForm
	if errs := c.ShouldBindJSON(&orderInfoForm); errs != nil {
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

	rsp, err := global.OrderClient.CreateOrder(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.OrderRequest{
			UserId:  int32(claims.ID),
			Name:    orderInfoForm.Name,
			Mobile:  orderInfoForm.Mobile,
			Address: orderInfoForm.Address,
			Post:    orderInfoForm.Post,
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}

	//支付宝url
	payUrl := NewAlipay().BuildPayUrl(rsp.OrderSn, rsp.PayAmount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rsp,
		"pay_url": payUrl,
		"message": "success",
	})
}

func (o *order) List(c *gin.Context) {
	e, blockErr := sentinel.Entry("order-flow-warmup-resource", sentinel.WithTrafficType(base.Inbound))
	if blockErr != nil {
		zap.S().Warn(blockErr.Error())
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code":    1,
			"message": "请求过于频繁,请稍后再试",
		})
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

	// 只能解决 parentContext 问题, tracer问题依然没解决
	// parentSpan, _ := c.Get("parentSpan")
	// opentracing.ContextWithSpan(context.Background(), parentSpan.(opentracing.Span))

	rsp, err := global.OrderClient.QueryOrderList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.OrderRequest{
			UserId: int32(claims.ID),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	e.Exit()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rsp,
		"message": "success",
	})
}

func (o *order) Detail(c *gin.Context) {
	id := c.Param("id")
	order_id, _ := strconv.Atoi(id)

	claims, err := getJwtInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	rsp, err := global.OrderClient.QueryOrderDetail(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.OrderRequest{
			UserId: int32(claims.ID),
			Id:     int32(order_id),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}

	var payUrl string
	if rsp.Status != "TRADE_SUCCESS" {
		//支付链接
		payUrl = NewAlipay().BuildPayUrl(rsp.OrderSn, rsp.PayAmount)
	}

	// c.Redirect(http.StatusTemporaryRedirect, payUrl)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rsp,
		"pay_url": payUrl,
		"message": "success",
	})
}

// delete order
func (o *order) Delete(c *gin.Context) {
	// orderSn := c.PostForm("order_sn")
	orderId := c.PostForm("order_id")
	order_id, _ := strconv.Atoi(orderId)

	claims, err := getJwtInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}
	_, err = global.OrderClient.DelOrder(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.OrderRequest{
			Id:     int32(order_id),
			UserId: int32(claims.ID),
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
