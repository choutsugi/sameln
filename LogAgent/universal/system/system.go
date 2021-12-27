package system

import (
	"LogAgent/universal/error"
	"net"
	"time"
)

// LocalTime 获取本地时间（毫秒级）
func LocalTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05.000")
}

func UtcTime() string {
	return time.Now().Local().Format("2006-01-02T15:04:05.000+0800")
}

func LocalIP() (ip string, err *error.Error) {
	address, raw := net.InterfaceAddrs()
	if raw != nil {
		return "", error.NewError(raw, error.CodeEtcdGetLocalIpFailed)
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
	return "", error.NewError(raw, error.CodeEtcdGetLocalIpFailed)
}
