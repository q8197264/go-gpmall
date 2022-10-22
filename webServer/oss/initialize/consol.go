package initialize

import (
	"fmt"
	"webServer/oss/global"
	"webServer/oss/utils/register/consul"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var serviceID string = uuid.NewV4().String()

func Register() {
	c := consul.NewClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port, serviceID)
	if err := c.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Consul.Tags,
	); err != nil {
		zap.S().Warnf("%s 服务注册失败: %s", global.ServerConfig.Name, err.Error())
		fmt.Printf("%s 服务注册失败: %s\n", global.ServerConfig.Name, err.Error())
	} else {
		zap.S().Infof("%s 服务注册成功: %s", global.ServerConfig.Name, serviceID)
		fmt.Printf("%s 服务注册成功: %s\n", global.ServerConfig.Name, serviceID)
	}
}

func Deregister() {
	c := consul.NewClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port, serviceID)
	if err := c.Deregister(); err != nil {
		zap.S().DPanicf("%s 服务注销失败: %s", global.ServerConfig.Name, err.Error())
		fmt.Printf("%s 服务注销失败: %s\n", global.ServerConfig.Name, err.Error())
	} else {
		zap.S().Infof("%s 服务注销: %s", global.ServerConfig.Name, serviceID)
		fmt.Printf("%s 服务注销成功: %s\n", global.ServerConfig.Name, serviceID)
	}
}
