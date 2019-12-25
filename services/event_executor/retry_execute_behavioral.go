package event_executor

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/pkg/rabbit"
	"time"
)

const errorTopic = "%s_ERROR"

type kafka2RabbitRetry struct {
	storageData    TopicExchangeData
	rabbitProducer rabbit.Producer
	kafkaProducer  sarama.SyncProducer
}

func RetryBehavioral(kafkaProducer sarama.SyncProducer, rabbitProducer rabbit.Producer, storageData TopicExchangeData) Executor {
	return &kafka2RabbitRetry{
		kafkaProducer:  kafkaProducer,
		rabbitProducer: rabbitProducer,
		storageData:    storageData,
	}
}

func (k *kafka2RabbitRetry) Execute(message *sarama.ConsumerMessage) (err error) {
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		err = execute(ctx, message, k.rabbitProducer, k.storageData)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		fmt.Printf("Message is not published to topic: %+v so is routing to error topic: %+v, message: %+v", message.Topic, fmt.Sprintf(errorTopic, message.Topic))
		_, _, err = k.sendToErrorTopic(message)
		if err != nil {
			fmt.Printf("Message is not published to error topic: %+v", fmt.Sprintf(errorTopic, message.Topic))
		}
	}
	return err
}

func (k *kafka2RabbitRetry) sendToErrorTopic(message *sarama.ConsumerMessage) (int32, int64, error) {
	return k.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: fmt.Sprintf(errorTopic, message.Topic),
		Value: sarama.StringEncoder(message.Value),
	})
}
