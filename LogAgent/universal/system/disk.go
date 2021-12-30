package system

import (
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"github.com/shirou/gopsutil/disk"
)

type UsageStat = disk.UsageStat

type DiskInfo struct {
	PartitionUsageStat map[string]*UsageStat
}

func GetDiskInfo() (info *DiskInfo, err *error.Error) {
	info = &DiskInfo{
		PartitionUsageStat: make(map[string]*UsageStat, 16),
	}
	stats, raw := disk.Partitions(true)
	if raw != nil {
		logger.L().Errorf("The System module gets disk-info unsuccessfully! Error:%s", raw.Error())
		return info, error.NewError(raw, codes.SystemGetDiskInfoFailed)
	}

	for _, stat := range stats {
		usage, raw := disk.Usage(stat.Mountpoint)
		if raw != nil {
			logger.L().Errorf("The System module gets disk-info unsuccessfully! Error:%s", raw.Error())
			continue
		}
		info.PartitionUsageStat[stat.Mountpoint] = usage
	}

	logger.L().Debug("The System module gets disk-info successfully!")
	return info, error.Null()
}
