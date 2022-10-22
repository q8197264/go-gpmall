package global

import (
	"webServer/goods/config"
	"webServer/goods/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans        ut.Translator
	GoodsClient  proto.GoodsClient
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
