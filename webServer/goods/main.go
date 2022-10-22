package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"webServer/goods/global"
	"webServer/goods/initialize"

	"go.uber.org/zap"
)

func main() {
	// 全局日志配置路径
	initialize.InitLogger()

	// 加载配置
	initialize.LoadConfig()

	// 注册当前服务
	initialize.ServiceRegister()

	// 加载grpc connect pool
	initialize.GrpcBalancer()

	// 初始化翻译
	if err := initialize.TranslateInit("zh"); err != nil {
		zap.S().Infof("初始化错误提示翻译失败")
	}

	// 路由配置
	router := initialize.InitRouter()

	go func() {
		address := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
		zap.S().Infof("web商品服务启动:%s", address)
		err := router.Run(address)
		if err != nil {
			zap.S().Fatal("err:", err.Error())
		} else {
			zap.S().Info("web商品服务启动:", address)
		}
	}()

	exit()
}

func exit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.ServiceDeregister()
	println("服务停止")
}
