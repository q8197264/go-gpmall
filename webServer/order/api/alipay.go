package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"webServer/order/global"
	"webServer/order/proto"
)

const serverDomain string = "http://ae57i5.natappfree.cc/v1/pay"

type payClient struct {
	appid        string
	privateKey   string
	publicKey    string
	serverDomain string
	callbackUrl  string
	notifyUrl    string
	alipayClient *alipay.Client
}

func NewAlipay() *payClient {
	client := &payClient{
		appid:       global.ServerConfig.Alipay.AppId,
		privateKey:  global.ServerConfig.Alipay.PrivateKey,
		publicKey:   global.ServerConfig.Alipay.AliPublicKey,
		callbackUrl: serverDomain + global.ServerConfig.Alipay.CallbackUrl,
		notifyUrl:   serverDomain + global.ServerConfig.Alipay.NotifyUrl,
	}

	var err error
	if client.alipayClient, err = alipay.New(client.appid, client.privateKey, false); err != nil {
		fmt.Printf("初始化支付宝失败:%s\n", err)
		return nil
	}

	// load alipay public key
	if err = client.alipayClient.LoadAliPayPublicKey(client.publicKey); err != nil {
		fmt.Printf("加载公钥发生错误:%s\n", err)
	}

	// 使用支付宝证书
	if err = client.alipayClient.LoadAppPublicCertFromFile("appCertPublicKey_2016073100129537.crt"); err != nil {
		fmt.Printf("加载证书发生错误:%s\n", err)
		// return nil
	}

	if err = client.alipayClient.LoadAliPayRootCertFromFile("alipayRootCert.crt"); err != nil {
		fmt.Printf("加载证书发生错误:%s\n", err)
		// return nil
	}
	if err = client.alipayClient.LoadAliPayPublicCertFromFile("alipayCertPublicKey_RSA2.crt"); err != nil {
		fmt.Printf("加载证书发生错误:%s\n", err)
		// return nil
	}

	return client
}

func (this *payClient) BuildPayUrl(orderSn string, totalAmount float32) string {
	var tradeNo = fmt.Sprintf("%s", orderSn)
	var payAmount = fmt.Sprintf("%.2f", totalAmount)

	var p = alipay.TradePagePay{}
	p.NotifyURL = this.notifyUrl
	p.ReturnURL = this.callbackUrl
	p.Subject = "支付测试:" + tradeNo
	p.OutTradeNo = tradeNo
	p.TotalAmount = payAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, _ := this.alipayClient.TradePagePay(p)
	// println(url.String())

	return url.String()
}

/*
	map[
		app_id:[2021000119660492]
		auth_app_id:[2021000119660492]
		charset:[utf-8]
		method:[alipay.trade.page.pay.return]
		out_trade_no:[202204151711381655]
		seller_id:[2088621958287916]
		sign:[r2Qmj433QmO27yOEkTOuIwvifTorj6GPULx8wH4HTj7zALpoyNKWdTA1hpU+fduKarPHED2Q45yBKh74JguVBVZPEaYGWKrlIr8pLkEpQWrMf5L2xR56ukrTB3nPXCxFihjHGUMbRCd5qbBuIC4tySd9o8WKJuwlO8N++vnKnjSX7UXl64btrmuTGdsqA/k5sVUc4ULCRr3DuXOJxvuvn0hqbqvNClkLC/1fC8lV96uIr9DadQSwcUzzojt48rlvGlrfCqXPafqlhIz2yHwbs0RLV91uTBREulveyhyP75/7gZfLBJmkQMntE5BWeeB1ylcDUbJLOAo1+cHtP8tbTg==]
		sign_type:[RSA2]
		timestamp:[2022-05-20 23:33:43]
		total_amount:[43.80]
		trade_no:[2022052022001426360501878833]
		version:[1.0]
	]
app_id=2021000119660492&auth_app_id=2021000119660492&charset=utf-8&method=alipay.trade.page.pay.return&out_trade_no=202204151711381655&seller_id=2088621958287916&sign=r2Qmj433QmO27yOEkTOuIwvifTorj6GPULx8wH4HTj7zALpoyNKWdTA1hpU+fduKarPHED2Q45yBKh74JguVBVZPEaYGWKrlIr8pLkEpQWrMf5L2xR56ukrTB3nPXCxFihjHGUMbRCd5qbBuIC4tySd9o8WKJuwlO8N++vnKnjSX7UXl64btrmuTGdsqA/k5sVUc4ULCRr3DuXOJxvuvn0hqbqvNClkLC/1fC8lV96uIr9DadQSwcUzzojt48rlvGlrfCqXPafqlhIz2yHwbs0RLV91uTBREulveyhyP75/7gZfLBJmkQMntE5BWeeB1ylcDUbJLOAo1+cHtP8tbTg==&sign_type=RSA2&timestamp=022-05-20 23:33:43&total_amount=43.80&trade_no=2022052022001426360501878833&version=1.0
*/
func (this *payClient) Callback(c *gin.Context) {
	c.Request.ParseForm()
	zap.S().Infof("%+v\n", c.Request.Form)
	//ok, err := aliClient.VerifySign(c.Request.Form)
	//if err != nil {
	//	zap.S().Infof("回调验证签名发生错误", err)
	//	return
	//}
	//
	//if ok == false {
	//	zap.S().Infof("回调验证签名未通过")
	//	return
	//}

	var outTradeNo = c.Request.Form.Get("out_trade_no")
	// total_amount
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := this.alipayClient.TradeQuery(p)
	if err != nil {
		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())
		return
	}
	if rsp.IsSuccess() == false {
		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
		return
	}

	c.String(http.StatusOK, "订单 %s 支付成功", outTradeNo)
}

/*
	map[
		app_id:[2021000119660492]
		auth_app_id:[2021000119660492]
		buyer_id:[2088622958326362]
		buyer_pay_amount:[43.80]
		charset:[utf-8]
		fund_bill_list:[
			[{"amount":"43.80","fundChannel":"ALIPAYACCOUNT"}]
		]
		gmt_create:[2022-05-20 23:33:00]
		gmt_payment:[2022-05-20 23:33:34]
		invoice_amount:[43.80]
		notify_id:[2022052000222233335026360519897668]
		notify_time:[2022-05-20 23:33:37]
		notify_type:[trade_status_sync]
		out_trade_no:[202204151711381655]
		point_amount:[0.00]
		receipt_amount:[43.80]
		seller_id:[2088621958287916]
		sign:[sAPvkAAzOWHXNtU2fkmPVJ82APRS04qr7yPvn3Fzo+uYXSgcFHVMgqQJp7gSRWVpFZOD0dRIcHTkfdlEdAdK1lkVTuDAu4uFVsozwSh9o/rOTpeL1wwR4Xn+9IzHqgsh6V7D5vj1p8J+P8yUH3wH4RBR8iDSiLZqNj84V+ENtHkrDut5V+UX3lbKbBUlg/WmoPkGnhWTxIbtH5xrjy4MUdmkIOZT+5f8WRkHgLCrgIw2Rz5lWzzYvvbRmtmpgPJ8Jw485bL/NnmLOLf9xLzKYOLXcyQtsjYquXv3zMAKxwn2XE/TLVjn7UdzaWT7RUAhdW0h5PWvi3kGdzoNyVSmCg==]
		sign_type:[RSA2]
		subject:[支付测试:202204151711381655]
		total_amount:[43.80]
		trade_no:[2022052022001426360501878833]
		trade_status:[TRADE_SUCCESS]
		version:[1.0]
	]
*/
func (this *payClient) Notify(c *gin.Context) {
	var noti, _ = this.alipayClient.GetTradeNotification(c.Request)
	if noti != nil {
		fmt.Println("交易状态为:", noti.TradeStatus)
	}
	fmt.Printf("-- %+v\n", noti)

	// 二、
	c.Request.ParseForm()
	ok, err := this.alipayClient.VerifySign(c.Request.Form)
	if err != nil {
		zap.S().Infof("异步通知验证签名发生错误", err.Error())
		return
	}

	if ok == false {
		zap.S().Infof("异步通知验证签名未通过")
		return
	}

	var outTradeNo = c.Request.Form.Get("out_trade_no")
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := this.alipayClient.TradeQuery(p)
	if err != nil {
		zap.S().Infof("异步通知验证订单 %s 信息发生错误: %s \n", outTradeNo, err.Error())
		return
	}
	if rsp.IsSuccess() == false {
		zap.S().Infof("异步通知验证订单 %s 信息发生错误: %s-%s \n", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
		return
	}

	// 变更订单状态
	_, err = global.OrderClient.UpdateOrderStatus(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.OrderStatusRequest{
			OrderSn: outTradeNo,
			Status:  "TRADE_SUCCESS",
		},
	)
	if err != nil {
		zap.S().Errorf("支付通知订单修改失败:", err.Error())
		return
	}
	zap.S().Infof("订单 %s 支付成功 \n", outTradeNo)

	c.String(http.StatusOK, "success")
}
