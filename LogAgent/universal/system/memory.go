package system

import (
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"github.com/shirou/gopsutil/mem"
)

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

func GetMemoryInfo() (info *MemoryInfo, err *error.Error) {
	info = new(MemoryInfo)
	stat, raw := mem.VirtualMemory()
	if raw != nil {
		logger.L().Errorf("The System module gets memory-info unsuccessfully! Error:%s", raw.Error())
		return info, error.NewError(raw, codes.SystemGetMemoryInfoFailed)
	}

	info.Total = stat.Total
	info.Available = stat.Available
	info.Used = stat.Used
	info.UsedPercent = stat.UsedPercent
	info.Buffers = stat.Buffers
	info.Cached = stat.Cached

	logger.L().Debug("The System module gets memory-info successfully!")
	return info, error.Null()
}
