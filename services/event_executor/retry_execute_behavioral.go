package event_executor

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/pkg/rabbit"
	"strings"
	"time"
)

const retryTopicPrefix = "RETRY"
const errorTopicPrefix = "ERROR"

type kafka2RabbitRetry struct {
	storageData    TopicExchangeData
	rabbitProducer rabbit.Producer
	kafkaProducer  sarama.SyncProducer
	retryInterval  int
	retryCount     int
}

func NewRetryExecutor(retryInterval, retryCount int, kafkaProducer sarama.SyncProducer, rabbitProducer rabbit.Producer, storageData TopicExchangeData) Executor {
	return &kafka2RabbitRetry{
		kafkaProducer:  kafkaProducer,
		rabbitProducer: rabbitProducer,
		storageData:    storageData,
		retryCount:     retryCount,
		retryInterval:  retryInterval,
	}
}

func (k *kafka2RabbitRetry) Execute(message *sarama.ConsumerMessage) (err error) {
	ctx := context.Background()
	for i := 0; i < k.retryCount; i++ {
		err = execute(ctx, message, k.rabbitProducer, k.storageData)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(k.retryInterval) * time.Second)
	}
	if err != nil {
		fmt.Printf("Message is not published to topic: %+v so is routing to error topic: %+v", message.Topic, fmt.Sprintf(errorTopicPrefix, message.Topic))
		if err := k.sendToErrorTopic(message); err != nil {
			fmt.Printf("Message is not published to error topic: %+v", fmt.Sprintf(errorTopicPrefix, message.Topic))
		}
	}
	return err
}

func (k *kafka2RabbitRetry) sendToErrorTopic(message *sarama.ConsumerMessage) error {
	errorTopic := strings.Replace(message.Topic, retryTopicPrefix, errorTopicPrefix, 0)
	_, _, err := k.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: fmt.Sprintf(errorTopic, message.Topic),
		Value: sarama.StringEncoder(message.Value),
	})
	return err
}
