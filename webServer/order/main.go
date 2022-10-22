package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"webServer/order/global"
	"webServer/order/initialize"
)

func main() {
	initialize.InitLogger()
	initialize.InitSentinel()
	initialize.InitNacosConfig()
	initialize.InitGrpcClient()
	initialize.ServiceRegister()
	initialize.InitTranslate("zh")
	engine := initialize.InitRouter()
	go func() {
		fmt.Printf("服务[%s]启动...%s:%d\n", global.ServerConfig.Name, global.ServerConfig.Host, global.ServerConfig.Port)
		engine.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	}()

	wait()
}

func wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initialize.ServiceDeregister()
	fmt.Printf("服务[%s]退出...\n", global.ServerConfig.Name)
}
