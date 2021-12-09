package blunder

const (
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

const (
	CodeSysSuccess = iota + 10000
	CodeSysSettingsInitSucceed
	CodeSysLoggerInitSucceed

	CodeSysUnknownAppMode
	CodeSysSettingsInitFailed
	CodeSysSettingsConfigUpdated
	CodeSysLoggerInitFailed
)

const (
	CodeEtcdConnectFailed = iota + 30000
	CodeEtcdGetLocalIpFailed
	CodeEtcdPutConfFailed
	CodeEtcdGetConfFailed
	CodeEtcdConfIsNotExist
	CodeEtcdConfUnmarshalFailed
	CodeEtcdConfUpdated
)

const (
	CodeKafkaConnFailed = iota + 40000
	CodeKafkaSendFailed
)

const (
	CodeNsqCreateProducerFailed = iota + 50000
	CodeNsqPublishFailed
)

const (
	CodeInfluxDbConnFailed = iota + 60000
	CodeInfluxDbCreatePointsFailed
	CodeInfluxDbCreatePointFailed
	CodeInfluxDbQueryFailed
	CodeInfluxDbWriteFailed
)

var blunderMsg = map[uint64]string{
	CodeSysSuccess:             "成功",
	CodeSysSettingsInitSucceed: "配置模块初始化成功",
	CodeSysLoggerInitSucceed:   "日志模块初始化成功",

	CodeSysSettingsInitFailed:    "配置模块初始化失败",
	CodeSysUnknownAppMode:        "未识别的应用启动模式，使用默认值：dev",
	CodeSysSettingsConfigUpdated: "配置文件已更新",
	CodeSysLoggerInitFailed:      "日志模块初始化失败（解析日志级别失败）",

	CodeEtcdConnectFailed:       "Etcd连接失败",
	CodeEtcdGetLocalIpFailed:    "Etcd获取本机IP失败",
	CodeEtcdPutConfFailed:       "Etcd设置数据失败",
	CodeEtcdGetConfFailed:       "Etcd获取数据失败",
	CodeEtcdConfIsNotExist:      "Etcd数据不存在",
	CodeEtcdConfUnmarshalFailed: "Etcd解析数据失败",
	CodeEtcdConfUpdated:         "Etcd数据已更新",

	CodeKafkaConnFailed: "Kafka连接失败",
	CodeKafkaSendFailed: "Kafka发送失败",

	CodeNsqCreateProducerFailed: "Nsq创建实例失败",
	CodeNsqPublishFailed:        "Nsq发布消息失败",

	CodeInfluxDbConnFailed:         "InfluxDB连接失败",
	CodeInfluxDbCreatePointsFailed: "InfluxDB批量创建数据点失败",
	CodeInfluxDbCreatePointFailed:  "InfluxDB创建数据点失败",
	CodeInfluxDbQueryFailed:        "InfluxDB查询失败",
	CodeInfluxDbWriteFailed:        "InfluxDB写入失败",
}

type Error struct {
	State string
	Code  uint64
	Msg   string
	Err   error
}
