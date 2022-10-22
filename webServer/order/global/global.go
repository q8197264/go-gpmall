package global

import (
	"webServer/order/config"
	"webServer/order/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfig    *config.NacosConfig
	Trans          ut.Translator
	OrderClient    proto.OrderClient
	ShopCartClient proto.ShopCartClient
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
)
