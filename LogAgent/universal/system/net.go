package system

import (
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"net"
	"strings"
)

func GetLocalIP() (ip string, err *error.Error) {
	address, raw := net.InterfaceAddrs()
	if raw != nil {
		return "", error.NewError(raw, codes.SystemGetLocalIpFailed)
	}

	for _, addr := range address {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		ip := ipAddr.IP.String()
		return ip, error.Null()
	}
	return "", error.NewError(raw, codes.SystemGetLocalIpFailed)
}

func GetLocalIpByDial() (ip string, err *error.Error) {
	conn, raw := net.Dial("udp", "8.8.8.8:80")
	if raw != nil {
		return
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(addr.IP.String(), ":")[0]
	return
}
