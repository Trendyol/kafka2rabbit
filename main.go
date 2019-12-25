package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	slack_api_client "github.com/kafka2rabbit/pkg/clients"
	"github.com/kafka2rabbit/pkg/kafka"
	"github.com/kafka2rabbit/pkg/rabbit"
	"github.com/kafka2rabbit/pkg/rest"
	"github.com/kafka2rabbit/services"
	"log"
)

var TopicConfigurations = []services.TopicExchangeData{
	{
		Topic:        "kafka topic",
		Exchange:     "rabbit exchange",
		ExchangeKind: "fanout",
		RoutingKey:   "routing",
	},
}

func main() {
	appPort := "8080"
	if err := initializeListeners(); err != nil {
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/_monitoring/ready", healthCheck)

	fmt.Println("[x] Kafka 2 Rabbit is running in " + appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Panicln(err.Error())
	}
}

func initializeListeners() error {
	for _, val := range TopicConfigurations {
		if err := runKafkaToRabbitListener(val); err != nil {
			fmt.Println("kafka rabbit listener has an error", err)
			return err
		}
	}
	return nil
}

func runKafkaToRabbitListener(data services.TopicExchangeData) error {
	config := kafka.ConnectionParameters{
		ConsumerGroupID: data.Topic + data.Exchange,
		ClientID:        "kafka2rabbit",
		Brokers:         "brokers",
		KafkaUsername:   "username",
		KafkaPassword:   "password",
		Version:         "version",
		Topic:           data.Topic,
		FromBeginning:   true,
	}

	connParameters := rabbit.ConnectionParameters{
		ConnectionString: "amqp://user:pass@addr:port",
		PrefetchCount:    15,
		RetryCount:       3,
		RetryInterval:    300,
	}

	rabbitProducer, err := rabbit.NewProducer(connParameters)
	if err != nil {
		return err
	}

	kafkaProducer, err := kafka.NewProducer(config)
	if err != nil {
		return err
	}

	if err := rabbitProducer.DeclareExchange(context.Background(), rabbit.Exchange{
		Kind: data.ExchangeKind,
		Name: data.Exchange,
	}); err != nil {
		return err
	}
	client := rest.NewClient()
	slackApiClient := slack_api_client.NewClient("slac_url",
		"slack username",
		"slack channel",
		client)
	eventHandler := services.NewEventHandler(rabbitProducer, data, kafkaProducer, slackApiClient)
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return err
	}
	consumer.Subscribe(eventHandler)
	fmt.Printf("%v listener is starting", data.Topic)
	return nil
}
func healthCheck(c *gin.Context) {
	c.JSON(200, "Healthy")
}
