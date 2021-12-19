package kafka

import (
	"LogAgent/error"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int64) *error.Error {
	// 1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	var err error.Raw
	// 2.连接kafka
	if client, err = sarama.NewSyncProducer(address, config); err != nil {
		zap.L().Fatal(error.GetInfo(error.CodeKafkaConnFailed))
		return error.NewError(err, error.CodeKafkaConnFailed)
	}

	// 3.初始化MsgChan
	msgChan = make(chan *sarama.ProducerMessage, chanSize)

	// 4.启动后台goroutine用于发送
	go sendMsg()

	return error.NullWithCode(error.CodeSysKafkaInitSucceed)
}

func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				zap.L().Warn(error.GetInfo(error.CodeKafkaSendFailed))
			}
			zap.L().Debug(
				error.GetInfo(error.CodeKafkaSendSucceed),
				zapcore.Field{Key: "pid", Interface: pid},
				zapcore.Field{Key: "offset", Interface: offset},
			)
		}
	}
}

func Write(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
