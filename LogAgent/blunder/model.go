package blunder

const (
	EtcdConnectFailed = iota + 30000
	EtcdGetLocalIpFailed
	EtcdPutConfFailed
	EtcdGetConfFailed
	EtcdConfIsNotExist
	EtcdConfUnmarshalFailed
	EtcdConfUpdated
)

const (
	KafkaConnFailed = iota + 40000
	KafkaSendFailed
)

const (
	NsqCreateProducerFailed = iota + 50000
	NsqPublishFailed
)

const (
	InfluxDbConnFailed = iota + 60000
	InfluxDbCreatePointsFailed
	InfluxDbCreatePointFailed
	InfluxDbQueryFailed
	InfluxDbWriteFailed
)

var blunderMsg = map[uint64]string{
	EtcdConnectFailed:       "Etcd连接失败",
	EtcdGetLocalIpFailed:    "Etcd获取本机IP失败",
	EtcdPutConfFailed:       "Etcd设置数据失败",
	EtcdGetConfFailed:       "Etcd获取数据失败",
	EtcdConfIsNotExist:      "Etcd数据不存在",
	EtcdConfUnmarshalFailed: "Etcd解析数据失败",
	EtcdConfUpdated:         "Etcd数据已更新",

	KafkaConnFailed: "Kafka连接失败",
	KafkaSendFailed: "Kafka发送失败",

	NsqCreateProducerFailed: "Nsq创建实例失败",
	NsqPublishFailed:        "Nsq发布消息失败",

	InfluxDbConnFailed:         "InfluxDB连接失败",
	InfluxDbCreatePointsFailed: "InfluxDB批量创建数据点失败",
	InfluxDbCreatePointFailed:  "InfluxDB创建数据点失败",
	InfluxDbQueryFailed:        "InfluxDB查询失败",
	InfluxDbWriteFailed:        "InfluxDB写入失败",
}

type Errors struct {
	Msg string
	Err *error
}
