package event_executor

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/pkg/rabbit"
	"github.com/kafka2rabbit/services"
	"github.com/streadway/amqp"
	"time"
)

const retryTopic = "%s_RETRY"

type kafka2RabbitNormal struct {
	storageData    services.TopicExchangeData
	rabbitProducer rabbit.Producer
	kafkaProducer  sarama.SyncProducer
}

func NormalBehavioral(kafkaProducer sarama.SyncProducer, rabbitProducer rabbit.Producer, data services.TopicExchangeData) Executor {
	return &kafka2RabbitNormal{
		storageData:    data,
		kafkaProducer:  kafkaProducer,
		rabbitProducer: rabbitProducer,
	}
}

func (k *kafka2RabbitNormal) Execute(message *sarama.ConsumerMessage) (err error) {
	ctx := context.Background()
	if err = execute(ctx, message, k.rabbitProducer, k.storageData); err != nil {
		_, _, err = k.sendToRetryTopic(message)
		if err != nil {
			fmt.Printf("Have an error occurred while publishing to retry topic: %+v , err:%+v", fmt.Sprintf(retryTopic, message.Topic), err)
		}
	} else {
		fmt.Printf("Message was published successfully, message: %+v", string(message.Value))
	}
	return err
}

func (k *kafka2RabbitNormal) sendToRetryTopic(message *sarama.ConsumerMessage) (int32, int64, error) {
	return k.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: fmt.Sprintf(retryTopic, message.Topic),
		Value: sarama.StringEncoder(message.Value),
	})
}

func execute(ctx context.Context, message *sarama.ConsumerMessage, publisher rabbit.Producer, storageData services.TopicExchangeData) error {
	brokerChannel, err := publisher.CreateConfirmedChannel(ctx)
	if err != nil {
		return fmt.Errorf("confirmation channel has an error , err:%v", err)
	}
	defer func() {
		_ = brokerChannel.Close(ctx)
	}()
	confirm := brokerChannel.NotifyPublish(make(chan amqp.Confirmation, 1))
	err = publisher.ProduceWithChannel(ctx, storageData.RoutingKey, storageData.Exchange, string(message.Value), 0, brokerChannel)
	if err != nil {
		return fmt.Errorf("message was not published to exchange , err:%v", err)
	}

	select {
	case confirmation := <-confirm:
		if !confirmation.Ack {
			return fmt.Errorf("not confirmed")
		}
	case <-time.After(3 * time.Second):
		return fmt.Errorf("confirmation timeout")
	}
	return nil
}
