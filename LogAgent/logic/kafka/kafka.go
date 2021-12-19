package kafka

import (
	"LogAgent/common/error"
	"LogAgent/common/logger"
	"github.com/Shopify/sarama"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize uint64) *error.Error {
	// 1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	var err error.RawErr
	// 2.连接kafka
	if client, err = sarama.NewSyncProducer(address, config); err != nil {
		logger.L().Fatalw(error.GetInfo(error.CodeKafkaConnFailed), "err", err.Error())
		return error.NewError(err, error.CodeKafkaConnFailed)
	}

	// 3.初始化MsgChan
	msgChan = make(chan *sarama.ProducerMessage, chanSize)

	// 4.启动后台goroutine用于发送
	go sendMsg()

	return error.Null()
}

func Write(msg *sarama.ProducerMessage) {
	msgChan <- msg
}

func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logger.L().Warn(error.GetInfo(error.CodeKafkaConnFailed))
			}
			logger.L().Debugw(
				error.GetInfo(error.CodeKafkaSendSucceed), "pid", pid, "offset", offset,
			)
		}
	}
}
