package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"webServer/userop/global"
	"webServer/userop/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.ServiceRegister()
	initialize.InitTranslator()
	initialize.InitGrpcClient()
	engine := initialize.InitRouter()
	go func() {
		engine.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	}()

	wait()
}

func wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.ServiceDeregister()
	println("退出userop服务")
}
