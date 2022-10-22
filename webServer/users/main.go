package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"webServer/users/global"
	"webServer/users/initialize"
)

func main() {
	// 加载日志
	initialize.Logger()

	// load config file
	initialize.InitConfig()

	// initialize.GrpcClient()
	initialize.GrpcBalancer()

	// 自动注册服务中心
	// serviceID := uuid.NewV4().String()
	// initialize.AutomatiRegister(serviceID)
	initialize.ServiceRegister()

	// 错误验证中文
	if err := initialize.TranslateInit("zh"); err != nil {
		zap.S().Infof("初始化错误提示翻译失败")
	}

	// load gin router
	router := initialize.Routes()

	// start serve
	go func() {
		host := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
		err := router.Run(host)
		if err != nil {
			zap.S().Panic("启动web服务失败:", err.Error())
		} else {
			zap.S().Debug(fmt.Sprintf("启动服务器 [ %s ]", host))
		}
	}()

	// comandline abort server
	exit()
}

func exit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	initialize.ServiceDeRegister()

	println("服务已关闭")
}
