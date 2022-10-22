package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Dail err", err)
		return
	}
	defer conn.Close()

	// 主动写数据给服务器
	conn.Write([]byte("来自客户端的连接..."))
}
