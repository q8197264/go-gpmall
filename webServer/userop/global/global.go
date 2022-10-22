package global

import (
	"webServer/userop/config"
	"webServer/userop/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans             ut.Translator
	GrpcFavClient     proto.FavoritesClient
	GrpcAddressClient proto.AddressClient
	GrpcPostClient    proto.PostClient
	NacosConf         *config.NacosConf
	ServerConfig      *config.ServerConf = &config.ServerConf{}
)
