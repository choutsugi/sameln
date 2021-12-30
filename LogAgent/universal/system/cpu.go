package system

import (
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

func GetCpuInfo() (info *CpuInfo, err *error.Error) {
	info = new(CpuInfo)
	percent, raw := cpu.Percent(time.Second, false)
	if raw != nil {
		logger.L().Errorf("The System module gets cpu-info unsuccessfully! Err:%s", raw.Error())
		return info, error.NewError(raw, codes.SystemGetCpuInfoFailed)
	}
	info.CpuPercent = percent[0]
	logger.L().Debug("The System module gets cpu-info successfully!")
	return info, error.Null()
}
