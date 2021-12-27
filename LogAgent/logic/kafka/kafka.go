package kafka

import (
	"LogAgent/universal/error"
	"LogAgent/universal/generic"
	"LogAgent/universal/logger"
	"LogAgent/universal/settings"
	"github.com/Shopify/sarama"
	"go.uber.org/atomic"
	"time"
)

type ProducerMessage = sarama.ProducerMessage

var (
	client      sarama.SyncProducer
	msgChan     chan *ProducerMessage
	initialized atomic.Bool
)

func Init(kafkaConfig *settings.KafkaConfigType) *error.Error {
	if initialized.Load() {
		return error.Null()
	}
	address := []string{kafkaConfig.Addr}
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
	msgChan = make(chan *sarama.ProducerMessage, kafkaConfig.ChanSize)

	// 4.启动后台goroutine用于发送
	go sendMsg()
	initialized.Store(true)
	return error.Null()
}

func Write(msg *ProducerMessage) {
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

func Close() {
	for tick := 0; tick < generic.TryCloseWithMaxTime; tick++ {
		if raw := client.Close(); raw == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// StringEncoder 同sarama库的StringEncoder
type StringEncoder string

func (s StringEncoder) Encode() ([]byte, error.RawErr) {
	return []byte(s), nil
}

func (s StringEncoder) Length() int {
	return len(s)
}
