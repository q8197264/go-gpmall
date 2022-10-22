package utils

import "net"

func GetFreePort(port int) int {
	if port > 0 {
		return port
	}
	addr, err := net.ResolveTCPAddr("tcp", "0")
	if err != nil {
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
	}

	return listen.Addr().(*net.TCPAddr).Port
}
