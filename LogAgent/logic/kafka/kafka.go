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

func Init(kafkaConfig *settings.KafkaConfigType) *error.Error {
	if initialized.Load() {
		logger.L().Error("The Kafka module unable to re-initialize!")
		return error.NewError(nil, codes.InitKafkaFailed)
	}
	address := []string{kafkaConfig.Addr}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	var raw error.RawErr
	if cli, raw = sarama.NewSyncProducer(address, config); raw != nil {
		logger.L().Errorf("The Kafka module creates sync-producer unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.InitKafkaFailed)
	}

	queue = make(chan *sarama.ProducerMessage, kafkaConfig.ChanSize)

	go sendMsg()

	initialized.Store(true)

	return error.Null()
}

func Write(msg *ProducerMessage) {
	queue <- msg
}

func sendMsg() {
	var raw error.RawErr
	var partition int32
	var offset int64

	for {
		select {
		case msg := <-queue:
			partition, offset, raw = cli.SendMessage(msg)
			if raw != nil {
				logger.L().Warnf("The Kafka module sends message to the Kafka service unsuccessfully! Error:%s", raw.Error())
			}
			logger.L().Debugf("The Kafka module sends message to the Kafka service successfully! Partition:%v Offset:%v.", partition, offset)
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
