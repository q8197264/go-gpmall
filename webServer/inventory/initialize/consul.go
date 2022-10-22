package initialize

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"webServer/inventory/global"
	"webServer/inventory/utils/consul"
)

var (
	serviceID = uuid.NewV4().String()
)

func ConsulRegister() {
	consulClient := consul.NewClient(
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
		serviceID,
	)
	if err := consulClient.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Consul.Tags,
	); err != nil {
		zap.S().DPanicf("%s 服务注册失败: %s", global.ServerConfig.Name, err.Error())
		fmt.Printf("%s 服务注册失败: %s", global.ServerConfig.Name, err.Error())
	} else {
		zap.S().Debugf("%s 服务注册成功: %s", global.ServerConfig.Name, serviceID)
		fmt.Printf("%s 服务注册成功: %s", global.ServerConfig.Name, serviceID)
	}
}

func ConsulDeregister() {
	consulClient := consul.NewClient(
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
		serviceID,
	)
	err := consulClient.Deregister(global.ServerConfig.Name, serviceID)
	if err != nil {
		zap.S().DPanicf("%s 服务注销失败: %s", global.ServerConfig.Name, err.Error())
		fmt.Printf("%s 服务注销失败: %s", global.ServerConfig.Name, err.Error())
	} else {
		zap.S().Debugf("%s 服务注销成功: %s", global.ServerConfig.Name, serviceID)
		fmt.Printf("%s 服务注销成功: %s", global.ServerConfig.Name, serviceID)
	}
}
