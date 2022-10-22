package global

import (
	ut "github.com/go-playground/universal-translator"

	"webServer/users/config"
	"webServer/users/proto"
)

var (
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
	UserClient   proto.UserClient
)
