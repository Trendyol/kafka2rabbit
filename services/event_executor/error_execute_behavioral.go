package event_executor

import (
	"fmt"
	"github.com/Shopify/sarama"
	slack_api_client "github.com/kafka2rabbit/pkg/clients"
)

const notificationFormat = "***************KAFKA_2_RABBIT*************** \n %v"

type kafka2Rabbit struct {
	client slack_api_client.Client
}

func (k *kafka2Rabbit) Execute(message *sarama.ConsumerMessage) error {
	err := k.client.PushNotification(fmt.Sprintf(notificationFormat, string(message.Value)))
	if err != nil {
		fmt.Println("Slack publish error, err:", err)
	}
	return nil
}

func NewErrorExecutor(client slack_api_client.Client) Executor {
	return &kafka2Rabbit{
		client: client,
	}
}
