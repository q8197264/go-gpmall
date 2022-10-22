package global

import (
	"webServer/inventory/config"
	"webServer/inventory/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans        *ut.Translator
	GrpcClient   proto.InventoryClient
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
