package util

import (
	"net"
	"strings"

	"logagent/common"
)

func GetOutBoundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", common.CannotGetLocalIPErr
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return strings.Split(localAddr.IP.String(), ":")[0], nil
}
