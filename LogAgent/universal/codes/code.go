// Package codes: status code
package codes

import "fmt"

type RawErr = error

const (
	Succeed             = iota + 10000
	InitSettingsSucceed = iota + 11000
	InitLoggerSucceed
	InitKafkaSucceed
	InitNsqSucceed
	InitEtcdSucceed
	InitInfluxDbSucceed

	KafkaConnectSucceed = iota + 12000
	KafkaSendSucceed

	NsqCreateProducerSucceed = iota + 13000
	NsqPublishSucceed

	EtcdConnectSucceed = iota + 14000
	EtcdGetIpSucceed
	EtcdSetConfigSucceed
	EtcdGetConfigSucceed
	EtcdConfigParseSucceed

	InfluxDbConnectSucceed = iota + 15000
	InfluxDbCreatePointsSucceed
	InfluxDbCreatePointSucceed
	InfluxDbQuerySucceed
	InfluxDbInsertSucceed
)

const (
	Failed             = iota + 20000
	InitSettingsFailed = iota + 21000
	InitLoggerFailed
	InitKafkaFailed
	InitNsqFailed
	InitEtcdFailed
	InitInfluxDbFailed

	KafkaConnectFailed = iota + 22000
	KafkaCreateProducerFailed
	KafkaSendFailed

	NsqCreateProducerFailed = iota + 23000
	NsqPublishFailed

	EtcdConnectFailed = iota + 24000
	EtcdGetIpFailed
	EtcdSetConfigFailed
	EtcdGetConfigFailed
	EtcdConfigNotFound
	EtcdConfigParseFailed

	InfluxDbConnectFailed = iota + 25000
	InfluxDbCreatePointsFailed
	InfluxDbCreatePointFailed
	InfluxDbQueryFailed
	InfluxDbInsertFailed

	SystemGetLocalIpFailed = iota + 26000

	CollectorInitTaskFailed = iota + 27000
)

const (
	Unknown = iota + 30000
)

var messages = map[uint64]string{
	Succeed:             "SUCCESS",
	InitSettingsSucceed: "Initialize the Settings module successfully.",
	InitLoggerSucceed:   "Initialize the Logger module successfully.",
	InitKafkaSucceed:    "Initialize the Kafka module successfully.",
	InitNsqSucceed:      "Initialize the Nsq module successfully.",
	InitEtcdSucceed:     "Initialize the Etcd module successfully.",
	InitInfluxDbSucceed: "Initialize the InfluxDb module successfully.",

	KafkaConnectSucceed: "The Kafka module connects to Kafka service successfully",
	KafkaSendSucceed:    "The Kafka module sends message to the Kafka service successfully",

	NsqCreateProducerSucceed: "The Nsq module creates producer successfully",
	NsqPublishSucceed:        "The Nsq module publishes message successfully",

	EtcdConnectSucceed:     "The Etcd module connects Etcd service successfully",
	EtcdGetIpSucceed:       "The Etcd module gets local ip successfully",
	EtcdSetConfigSucceed:   "The Etcd module sets config successfully",
	EtcdGetConfigSucceed:   "The Etcd module gets config successfully",
	EtcdConfigParseSucceed: "The Etcd module parses config successfully",

	InfluxDbConnectSucceed:      "The InfluxDb module connects InfluxDb service successfully",
	InfluxDbCreatePointsSucceed: "The InfluxDb module creates points successfully",
	InfluxDbCreatePointSucceed:  "The InfluxDb module creates point successfully",
	InfluxDbQuerySucceed:        "The InfluxDb module query data successfully",
	InfluxDbInsertSucceed:       "The InfluxDb module inserts data successfully",

	Failed:             "FAIL",
	InitSettingsFailed: "Initialize the Settings module unsuccessfully!",
	InitLoggerFailed:   "Initialize the Logger module unsuccessfully!",
	InitKafkaFailed:    "Initialize the Kafka module unsuccessfully!",
	InitNsqFailed:      "Initialize the Nsq module unsuccessfully!",
	InitEtcdFailed:     "Initialize the Etcd module unsuccessfully!",
	InitInfluxDbFailed: "Initialize the InfluxDb module unsuccessfully!",

	KafkaConnectFailed:        "The Kafka module connects to Kafka service unsuccessfully!",
	KafkaCreateProducerFailed: "The Kafka module creates sync-producer unsuccessfully!",
	KafkaSendFailed:           "The Kafka module sends message to the Kafka service unsuccessfully!",

	NsqCreateProducerFailed: "The Nsq module creates producer unsuccessfully!",
	NsqPublishFailed:        "The Nsq module publishes message unsuccessfully!",

	EtcdConnectFailed:     "The Etcd module connects Etcd service unsuccessfully!",
	EtcdGetIpFailed:       "The Etcd module gets local ip unsuccessfully!",
	EtcdSetConfigFailed:   "The Etcd module sets config unsuccessfully!",
	EtcdGetConfigFailed:   "The Etcd module gets config unsuccessfully!",
	EtcdConfigNotFound:    "The Etcd module cannot find the specified configuration from the Etcd service!",
	EtcdConfigParseFailed: "The Etcd module parses config unsuccessfully!",

	InfluxDbConnectFailed:      "The InfluxDb module connects InfluxDb service unsuccessfully!",
	InfluxDbCreatePointsFailed: "The InfluxDb module creates points unsuccessfully!",
	InfluxDbCreatePointFailed:  "The InfluxDb module creates point unsuccessfully!",
	InfluxDbQueryFailed:        "The InfluxDb module query data unsuccessfully!",
	InfluxDbInsertFailed:       "The InfluxDb module inserts data unsuccessfully!",

	SystemGetLocalIpFailed: "The System module gets local ip unsuccessfully!",

	CollectorInitTaskFailed: "The Collector module initializes tail-task unsuccessfully!",

	Unknown: "Unknown!",
}

func Message(code uint64) string {
	msg, isExist := messages[code]
	if !isExist {
		msg = "Unknown!"
	}
	return msg
}

func MessageWithRaw(code uint64, raw RawErr) string {
	if raw == nil {
		return Message(code)
	}
	return fmt.Sprintf(Message(code)+" Err:%s", raw.Error())
}
