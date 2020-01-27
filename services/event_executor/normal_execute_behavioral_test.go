package event_executor_test

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/pkg/kafka"
	"github.com/kafka2rabbit/pkg/rabbit"
	. "github.com/kafka2rabbit/pkg/util"
	"github.com/kafka2rabbit/services/event_executor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NormalExecuteBehavioral", func() {
	Context("with no error", func() {
		var (
			message     []byte
			topic       = "normal"
			exchange    = "exchange"
			queueName   = "queue"
			routingKey  = "routingKey"
			storageData = event_executor.TopicExchangeData{
				Topic:        topic,
				Exchange:     exchange,
				ExchangeKind: "Topic",
				RoutingKey:   routingKey,
			}
			rabbitApi RabbitApi
		)
		BeforeAll(func() {
			rabbitApi = NewRabbitApi(rabbitContainer)
			producerConf := kafka.ConnectionParameters{
				Version:         "2.2.0",
				ConsumerGroupID: "some-id",
				ClientID:        "oms-event-generator",
				Topic:           topic,
				Brokers:         kafkaContainer.Address(),
			}
			connParameters := rabbit.ConnectionParameters{
				ConnectionString: rabbitContainer.Address(),
				PrefetchCount:    1,
				RetryCount:       3,
				RetryInterval:    300,
			}
			publisher, err := rabbit.NewProducer(connParameters)
			Expect(err).NotTo(HaveOccurred())
			err = publisher.DeclareExchange(context.Background(), rabbit.Exchange{
				Kind: "topic",
				Name: exchange,
			})
			Expect(err).NotTo(HaveOccurred())
			err = rabbitApi.CreateQueue(queueName, exchange, routingKey)
			producer, err := kafka.NewProducer(producerConf)
			Expect(err).NotTo(HaveOccurred())

			executor := event_executor.NewEventExecutor()
			executor.SetStrategy(event_executor.NewNormalExecutor(producer, publisher, storageData))
			msg := &sarama.ConsumerMessage{
				Value: []byte("message"),
			}
			err = executor.Execute(msg)
			Expect(err).NotTo(HaveOccurred())

			message = rabbitApi.GetMessage(queueName)
		})

		AfterAll(func() {
			rabbitApi.PurgeQueue(queueName)
		})

		It("should consumed events count to be one", func() {
			Expect(string(message)).Should(Equal("message"))
		})
	})
	Context("with error when event publishing to rabbit", func() {
		var (
			message         []byte
			topic           = "normal"
			expectedPayload = "message"
			exchange        = "exchangeFailed"
			queueName       = "queue2"
			routingKey      = "routing2"
			storageData     = event_executor.TopicExchangeData{
				Topic:        topic,
				Exchange:     exchange,
				ExchangeKind: "topic",
				RoutingKey:   routingKey,
			}
			rabbitApi RabbitApi
		)

		BeforeAll(func() {
			rabbitApi = NewRabbitApi(rabbitContainer)
			producerConf := kafka.ConnectionParameters{
				Version:         "2.2.0",
				ConsumerGroupID: "some-id",
				ClientID:        "oms-event-generator",
				Topic:           topic + "_RETRY",
				Brokers:         kafkaContainer.Address(),
			}
			connParameters := rabbit.ConnectionParameters{
				ConnectionString: rabbitContainer.Address(),
				PrefetchCount:    1,
				RetryCount:       3,
				RetryInterval:    300,
			}
			producer, err := kafka.NewProducer(producerConf)
			Expect(err).NotTo(HaveOccurred())

			publisher, err := rabbit.NewProducer(connParameters)
			Expect(err).NotTo(HaveOccurred())

			Expect(err).NotTo(HaveOccurred())
			err = publisher.DeclareExchange(context.Background(), rabbit.Exchange{
				Kind: "topic",
				Name: exchange,
			})
			Expect(err).NotTo(HaveOccurred())
			err = rabbitApi.CreateQueue(queueName, exchange, routingKey)

			executor := event_executor.NewEventExecutor()
			executor.SetStrategy(event_executor.NewNormalExecutor(producer, publisher, storageData))
			msg := &sarama.ConsumerMessage{
				Topic: topic,
				Value: []byte("message"),
			}
			err = executor.Execute(msg)
			Expect(err).NotTo(HaveOccurred())

			message = rabbitApi.GetMessage(queueName)

		})
		It("should consumed message value equals to expected payload from retryTopicPrefix topic", func() {
			Expect(string(message)).Should(Equal(expectedPayload))
		})
	})
})
