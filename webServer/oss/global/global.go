package global

import (
	"webServer/oss/config"

	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc"
)

var (
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
	ClientConn   *grpc.ClientConn
)
