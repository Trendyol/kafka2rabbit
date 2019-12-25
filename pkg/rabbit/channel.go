package rabbit

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitChannel interface {
	DeclareExchange(exchange Exchange) error
	PrefetchCount() int
	Confirmation(noWait bool) error
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	PublishMessage(context.Context, string, string, amqp.Publishing) error
	Close(context.Context) error
}

const (
	Topic  = "topic"
	Fanout = "fanout"
)

type Exchange struct {
	Kind string
	Name string
}

type rabbitChannel struct {
	channel       *amqp.Channel
	prefetchCount int
}

func newChannel(ctx context.Context, connection *amqp.Connection, prefetchCount int) (brokerChannel RabbitChannel, err error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel creation error,err:%v", err)
	}

	channel.Qos(prefetchCount, 0, false)

	return &rabbitChannel{
		prefetchCount: prefetchCount,
		channel:       channel,
	}, nil
}

func (c *rabbitChannel) Confirmation(noWait bool) error {
	return c.channel.Confirm(noWait)
}

func (c *rabbitChannel) NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	return c.channel.NotifyPublish(confirm)
}

func (c *rabbitChannel) DeclareExchange(exchange Exchange) error {
	err := c.channel.ExchangeDeclare(
		exchange.Name, // name of the exchange
		exchange.Kind, // type
		true,          // durable
		false,         // delete when complete
		false,         // internal
		false,         // noWait
		nil,           // arguments
	)

	if err != nil {
		return fmt.Errorf("exchange declaration error,err:%v", err)
	}
	return nil
}

func (c *rabbitChannel) PrefetchCount() int {
	return c.prefetchCount
}

func (c *rabbitChannel) PublishMessage(ctx context.Context, exchangeName string, routingKey string, data amqp.Publishing) error {
	return c.channel.Publish(exchangeName, routingKey, false, false, data)
}

func (c *rabbitChannel) Close(ctx context.Context) error {

	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("rabbitChannel close error,err:%v", err)
	}
	return nil
}
