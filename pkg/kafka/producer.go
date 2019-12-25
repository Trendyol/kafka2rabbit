package kafka

import (
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

func NewProducer(connectionParams ConnectionParameters) (sarama.SyncProducer, error) {
	v, err := sarama.ParseKafkaVersion(connectionParams.Version)
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = v
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 5                    // Retry up to 5 times to produce the message. Default is 3
	config.Producer.Retry.Backoff = 10 * time.Second
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Metadata.Retry.Max = 1
	config.Metadata.Retry.Backoff = 5 * time.Second

	syncProducer, err := sarama.NewSyncProducer(strings.Split(connectionParams.Brokers, ","), config)
	if err != nil {
		return nil, err
	}

	return syncProducer, nil
}