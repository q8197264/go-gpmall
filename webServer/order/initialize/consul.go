package initialize

import (
	"fmt"
	"webServer/order/global"
	cs "webServer/order/utils/register/consul"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var serviceId = uuid.NewV4().String()

func client() *cs.ConsulClient {
	client := cs.New(
		global.ServerConfig.Name,
		global.ServerConfig.Consul.Host,
		global.ServerConfig.Consul.Port,
		global.ServerConfig.Consul.Tags,
		serviceId,
	)
	return client
}

func ServiceRegister() {
	if err := client().ServiceRegister(global.ServerConfig.Host, global.ServerConfig.Port); err != nil {
		zap.S().DPanic(err.Error())
	}
	fmt.Printf("注册订单服务[%s]成功...\n", global.ServerConfig.Name)
}

func ServiceDeregister() {
	if err := client().ServiceDeregister(); err != nil {
		zap.S().DPanic(err.Error())
	}
	fmt.Printf("注销订单服务[%s]成功...\n", global.ServerConfig.Name)
}
