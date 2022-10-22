package utils

import "net"

func GetFreePort(port int) int {
	if port > 0 {
		return port
	}
	addr, err := net.ResolveTCPAddr("tcp", "0")
	if err != nil {
		panic(err.Error())
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	defer listen.Close()

	return listen.Addr().(*net.TCPAddr).Port
}
