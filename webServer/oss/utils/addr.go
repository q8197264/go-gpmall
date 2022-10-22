package utils

import (
	"net"

	"go.uber.org/zap"
)

func GetFreePort(defaultPort int) int {
	if defaultPort > 0 {
		return defaultPort
	}
	laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		zap.S().DPanic(err.Error())
	}

	listen, err := net.ListenTCP("tcp", laddr)
	if err != nil {
	}

	defer listen.Close()

	return listen.Addr().(*net.TCPAddr).Port
}
