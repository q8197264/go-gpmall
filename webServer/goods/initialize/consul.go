package initialize

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"webServer/goods/global"
	"webServer/goods/utils/register/consul"
)

var serviceID string = uuid.NewV4().String()

func ServiceRegister() {
	registerClient := consul.NewRegistryClient(
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
	)
	err := registerClient.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Consul.Tags,
		serviceID,
	)
	if err != nil {
		zap.S().Warnf("%s服务注册失败:%s", global.ServerConfig.Name, err.Error())
		fmt.Printf("%s服务注册失败:%s\n", global.ServerConfig.Name, err.Error())
	} else {
		zap.S().Infof("%s服务注册成功", global.ServerConfig.Name)
		fmt.Printf("%s服务注册success\n", global.ServerConfig.Name)
	}
}

func ServiceDeregister() {
	registerClient := consul.NewRegistryClient(
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
	)
	err := registerClient.DeRegister(serviceID)

	if err != nil {
		zap.S().Warnf("%s服务注销失败:%s", global.ServerConfig.Name, err.Error())
		fmt.Printf("%s服务注销失败:%s\n", global.ServerConfig.Name, err.Error())
	} else {
		zap.S().Infof("%s服务注销success", global.ServerConfig.Name)
		fmt.Printf("%s服务注销success\n", global.ServerConfig.Name)
	}
}
