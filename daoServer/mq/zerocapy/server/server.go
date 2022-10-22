package main

import (
	"fmt"
	"net"
)

func main() {
	// 指定服务器通信协议ip地址和端口号
	listener, err := net.Listen("tcp", "127.0.0.1:10000") // 这里的listen并不是监听，而是创建一个用于监听的socket
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务器等待客户端建立连接.......")
	// 阻塞监听客户端连接请求
	conn, err := listener.Accept() // 监听
	if err != nil {
		fmt.Println("accept err", err)
		return
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf) // n是字节数
	if err != nil {
		fmt.Println("conn Read err", err)
		return
	}
	fmt.Println("服务器收到的数据是：", string(buf[:n]))
	fmt.Println("服务器与客户端成功建立连接")
	defer conn.Close()
	// 读取客户端发送数据
}
