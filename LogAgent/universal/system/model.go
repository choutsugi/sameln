package system

const (
	CpuInfoType  = "cpu"
	MemInfoType  = "mem"
	DiskInfoType = "disk"
	NetInfoType  = "net"
)

type SysInfo struct {
	IP   string
	Type string
	Data interface{}
}
