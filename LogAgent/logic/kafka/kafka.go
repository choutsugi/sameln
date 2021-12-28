package kafka

import (
	"LogAgent/universal/codes"
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
	cli         sarama.SyncProducer
	queue       chan *ProducerMessage
	initialized atomic.Bool
)

func Init(kafkaConfig *settings.KafkaConfigType) (err *error.Error) {
	if initialized.Load() {
		return
	}
	address := []string{kafkaConfig.Addr}
	// 1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 2.连接kafka
	var raw error.RawErr
	if cli, raw = sarama.NewSyncProducer(address, config); raw != nil {
		err = error.NewError(raw, codes.KafkaConnectFailed)
		return
	}

	// 3.初始化MsgChan
	queue = make(chan *sarama.ProducerMessage, kafkaConfig.ChanSize)

	// 4.启动后台goroutine用于发送
	go sendMsg()
	initialized.Store(true)

	return error.Null()
}

func Write(msg *ProducerMessage) {
	queue <- msg
}

func sendMsg() {
	var err *error.Error
	var raw error.RawErr
	var partition int32
	var offset int64

	for {
		select {
		case msg := <-queue:
			partition, offset, raw = cli.SendMessage(msg)
			if err != nil {
				err = error.NewError(raw, codes.KafkaSendFailed)
				logger.L().Warn(err.Info())
			}
			logger.L().Debugf(codes.Message(codes.KafkaSendSucceed)+" partition:%v offset:%v.", partition, offset)
		}
	}
}

func Close() {
	for tick := 0; tick < generic.TryCloseWithMaxTime; tick++ {
		if raw := cli.Close(); raw == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

type StringEncoder string

func (s StringEncoder) Encode() ([]byte, error.RawErr) {
	return []byte(s), nil
}

func (s StringEncoder) Length() int {
	return len(s)
}
