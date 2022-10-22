package utils

import "net"

//  获取空闲端口 [适用于线上环境，本地测试建议用固定端口便于测试]
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listen.Close()

	port := listen.Addr().(*net.TCPAddr).Port

	return port, nil
}
