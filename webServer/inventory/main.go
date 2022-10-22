package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"webServer/inventory/global"
	"webServer/inventory/initialize"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitNacos()
	initialize.ConsulRegister()
	initialize.InitGrpcClient()
	initialize.InitTranslator("zh")
	router := initialize.InitRouter()
	go func() {
		err := router.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
		if err != nil {
			zap.S().DPanic(err.Error())
		}
	}()

	exit()
}

func exit() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.ConsulDeregister()
	println("退出服务")
}
