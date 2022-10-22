package initialize

import (
	"fmt"
	"webServer/users/global"
	"webServer/users/utils/register/consul"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var serviceID = uuid.NewV4().String()

func ServiceRegister() {
	client := consul.NewRegistryClient(
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
	)
	err := client.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Consul.Tags,
		serviceID,
	)
	if err != nil {
		zap.S().Warnf("%s服务注册失败: %s", err.Error())
	} else {
		zap.S().Infof("%s服务注册成功", global.ServerConfig.Name)
		fmt.Printf("%s服务注册成功\n", global.ServerConfig.Name)
	}
}

func ServiceDeRegister() {
	client := consul.NewRegistryClient(
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
	)
	err := client.DeRegister(serviceID)
	if err != nil {
		zap.S().Warnf("%s服务注销失败: %s", err.Error())
	} else {
		zap.S().Infof("%s服务注销成功", global.ServerConfig.Name)
		fmt.Printf("%s服务注销成功\n", global.ServerConfig.Name)
	}
}
