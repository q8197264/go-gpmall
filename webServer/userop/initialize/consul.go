package initialize

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"webServer/userop/global"
	"webServer/userop/utils/register/consul"
)

var serviceId = uuid.NewV4().String()

func client() *consul.ConsulClient {
	c := consul.NewClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port, serviceId)
	return c
}

func ServiceRegister() {
	err := client().Register(global.ServerConfig.Name, global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Consul.Tags)
	if err != nil {
		zap.S().Warnf("[%s]服务注册失败...%s", global.ServerConfig.Name, err.Error())
	} else {
		fmt.Printf("[%s]服务注册成功...\n", global.ServerConfig.Name)
	}
}

func ServiceDeregister() {
	if err := client().Deregister(); err != nil {
		zap.S().Warnf("[%s]服务注销失败... %s", global.ServerConfig.Name, err.Error())
	} else {
		fmt.Printf("[%s]服务注销成功...\n", global.ServerConfig.Name)
	}
}
