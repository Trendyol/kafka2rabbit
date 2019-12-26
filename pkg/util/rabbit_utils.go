package util

import (
	rabbitCotnainer "github.com/kafka2rabbit/internal/integrationcontainers/rabbit"
	"github.com/kafka2rabbit/pkg/rabbit"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"
)

type rabbitApi struct {
	rabbitContainer *rabbitCotnainer.Container
}

type RabbitApi interface {
	CreateQueue(name string, defaultExchange string, routingKey string) error
	GetMessage(queueName string) []byte
	PurgeQueue(name string)
}

func NewRabbitApi(rabbitMQContainer *rabbitCotnainer.Container) RabbitApi {
	return &rabbitApi{
		rabbitMQContainer,
	}
}

func (r *rabbitApi) CreateQueue(name string, defaultExchange string, routingKey string) error {
	connParameters := rabbit.ConnectionParameters{
		ConnectionString: r.rabbitContainer.Address(),
		PrefetchCount:    1,
		RetryCount:       3,
		RetryInterval:    300,
	}

	conn, err := amqp.Dial(connParameters.ConnectionString)
	Expect(err).NotTo(HaveOccurred())

	ch, err := conn.Channel()
	Expect(err).NotTo(HaveOccurred())
	_, err = ch.QueueDeclare(name,
		false,
		false,
		false,
		false,
		nil,
	)
	Expect(err).NotTo(HaveOccurred())
	err = ch.QueueBind(
		name,
		routingKey,
		defaultExchange,
		false,
		nil,
	)
	Expect(err).NotTo(HaveOccurred())

	return nil
}

func (r *rabbitApi) PurgeQueue(name string) {
	connParameters := rabbit.ConnectionParameters{
		ConnectionString: r.rabbitContainer.Address(),
		PrefetchCount:    1,
		RetryCount:       3,
		RetryInterval:    300,
	}

	conn, err := amqp.Dial(connParameters.ConnectionString)
	Expect(err).NotTo(HaveOccurred())

	ch, err := conn.Channel()
	Expect(err).NotTo(HaveOccurred())

	_, err = ch.QueuePurge(name, false)
	Expect(err).NotTo(HaveOccurred())
}

func (r *rabbitApi) GetMessage(queueName string) []byte {
	connParameters := rabbit.ConnectionParameters{
		ConnectionString: r.rabbitContainer.Address(),
		PrefetchCount:    1,
		RetryCount:       3,
		RetryInterval:    300,
	}

	conn, err := amqp.Dial(connParameters.ConnectionString)
	Expect(err).NotTo(HaveOccurred())

	ch, err := conn.Channel()
	Expect(err).NotTo(HaveOccurred())

	msg, _, err := ch.Get(queueName, false)
	Expect(err).NotTo(HaveOccurred())

	return msg.Body
}
