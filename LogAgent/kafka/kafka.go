package kafka

import (
	"github.com/Shopify/sarama"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int64) {

}
