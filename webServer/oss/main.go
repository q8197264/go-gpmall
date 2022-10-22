package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"webServer/oss/global"
	"webServer/oss/initialize"

	"go.uber.org/zap"
)

func main() {
	// 日志初始化
	initialize.InitLogger()
	// 配置中心初始化
	initialize.InitNacos()
	// 翻译初始化
	initialize.InitTranslator("zh")
	// grpc客户端连接
	initialize.InitGrpcClient()

	// 注册中心初始化
	initialize.Register()

	// 路由初始化
	router := initialize.InitRouter()
	go func() {
		if err := router.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)); err != nil {
			println(err.Error())
			zap.S().DPanic(err.Error())
		}
	}()

	exit()
}

func exit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.Deregister()
	println("服务停止")
}
