package kafka

import "C"
import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

type ConnectionParameters struct {
	Brokers         string
	ClientID        string
	Version         string
	Topic           string
	ConsumerGroupID string
	KafkaUsername   string
	KafkaPassword   string
	ApplicationPort string
	FromBeginning   bool
}

type Consumer interface {
	Subscribe(handler EventHandler)
	Unsubscribe()
}

type EventHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type kafkaConsumer struct {
	topic         string
	retryTopic    string
	errorTopic    string
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(connectionParams ConnectionParameters) (Consumer, error) {
	v, err := sarama.ParseKafkaVersion(connectionParams.Version)
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = v
	config.Consumer.Return.Errors = true
	config.ClientID = connectionParams.ClientID
	if connectionParams.FromBeginning {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.Metadata.Retry.Max = 1
	config.Metadata.Retry.Backoff = 5 * time.Second

	cg, err := sarama.NewConsumerGroup(strings.Split(connectionParams.Brokers, ","), connectionParams.ConsumerGroupID, config)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		topic:         connectionParams.Topic,
		retryTopic:    fmt.Sprintf("%s_RETRY", connectionParams.Topic),
		errorTopic:    fmt.Sprintf("%s_ERROR", connectionParams.Topic),
		consumerGroup: cg,
	}, nil
}

func (c *kafkaConsumer) Subscribe(handler EventHandler) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if err := c.consumerGroup.Consume(ctx, []string{c.topic, c.retryTopic, c.errorTopic}, handler);
				err != nil {
				panic(err)
			}

			if ctx.Err() != nil {
				fmt.Println(ctx.Err())
				return
			}
		}
	}()
	go func() {
		for err := range c.consumerGroup.Errors() {
			fmt.Println(err.Error())
			cancel()
		}
	}()
}

func (c *kafkaConsumer) Unsubscribe() {
	if err := c.consumerGroup.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
