package influx

import (
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"LogAgent/universal/settings"
	"LogAgent/universal/system"
	client "github.com/influxdata/influxdb1-client/v2"
	"go.uber.org/atomic"
	"time"
)

var (
	cli         client.Client
	config      client.BatchPointsConfig
	initialized atomic.Bool
	activated   atomic.Bool
)

func Init(influxConfig *settings.InfluxDbConfigType) (err *error.Error) {
	if initialized.Load() {
		logger.L().Errorf("The InfluxDb module unable to re-initialize!")
		return error.NewError(nil, codes.InitInfluxDbFailed)
	}

	if influxConfig.Active {
		activated.Store(true)
	}

	var raw error.RawErr
	cli, raw = client.NewHTTPClient(client.HTTPConfig{
		Addr: influxConfig.Addr,
	})

	if raw != nil {
		logger.L().Errorf("The InfluxDb module connects InfluxDb service unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbConnectFailed)
	}

	config = client.BatchPointsConfig{
		Precision: influxConfig.Precision,
		Database:  influxConfig.Database,
	}

	initialized.Store(true)
	return error.Null()
}

func Query(cmd string) (ret []client.Result, err *error.Error) {
	query := client.Query{
		Command:  cmd,
		Database: "log-server",
	}
	if rsp, raw := cli.Query(query); raw == nil {
		if rsp.Error() != nil {
			return ret, error.NewError(rsp.Error(), codes.InfluxDbQueryFailed)
		}
		ret = rsp.Results
	} else {
		return ret, error.NewError(raw, codes.InfluxDbQueryFailed)
	}
	return ret, error.Null()
}

func InsertCpuInfo(info *system.CpuInfo) *error.Error {
	if !activated.Load() {
		return error.Null()
	}
	points, raw := client.NewBatchPoints(config)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": info.CpuPercent,
	}

	point, raw := client.NewPoint("cpu", tags, fields, time.Now())
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}
	points.AddPoint(point)
	raw = cli.Write(points)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	logger.L().Info("The InfluxDb module inserts data successfully!")
	return error.Null()
}

func InsertMemInfo(info *system.MemoryInfo) *error.Error {
	if !activated.Load() {
		return error.Null()
	}
	points, raw := client.NewBatchPoints(config)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"total":        info.Total,
		"available":    info.Available,
		"used":         info.Used,
		"used_percent": info.UsedPercent,
		"buffers":      info.Buffers,
		"cached":       info.Cached,
	}

	point, raw := client.NewPoint("memory", tags, fields, time.Now())
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}
	points.AddPoint(point)
	raw = cli.Write(points)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	logger.L().Info("The InfluxDb module inserts data successfully!")
	return error.Null()
}

func InsertDiskInfo(info *system.DiskInfo) *error.Error {
	if !activated.Load() {
		return error.Null()
	}
	points, raw := client.NewBatchPoints(config)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	for k, v := range info.PartitionUsageStat {
		tags := map[string]string{"path": k}
		fields := map[string]interface{}{
			"total":               v.Total,
			"free":                v.Free,
			"used":                v.Used,
			"user_percent":        v.UsedPercent,
			"inodes_total":        v.InodesTotal,
			"inodes_used":         v.InodesUsed,
			"inodes_free":         v.InodesFree,
			"inodes_used_percent": v.InodesUsedPercent,
		}
		point, raw := client.NewPoint("disk", tags, fields, time.Now())
		if raw != nil {
			logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
			continue
		}
		points.AddPoint(point)
	}

	raw = cli.Write(points)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	logger.L().Info("The InfluxDb module inserts data successfully!")
	return error.Null()
}

func InsertNetInfo() (err *error.Error) {
	if !activated.Load() {
		return error.Null()
	}
	points, raw := client.NewBatchPoints(config)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	tags := map[string]string{"net": "net"}
	fields := map[string]interface{}{}

	point, raw := client.NewPoint("net", tags, fields, time.Now())
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}
	points.AddPoint(point)
	raw = cli.Write(points)
	if raw != nil {
		logger.L().Errorf("The InfluxDb module inserts data unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InfluxDbInsertFailed)
	}

	logger.L().Info("The InfluxDb module inserts data successfully!")
	return error.Null()
}
