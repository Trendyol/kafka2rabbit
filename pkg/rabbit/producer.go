package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type producer struct {
	client RabbitClient
}

type Producer interface {
	CreateConfirmedChannel(ctx context.Context) (RabbitChannel, error)
	ProduceWithChannel(ctx context.Context, routingKey, exchangeName string, payload interface{}, priority uint8, channel RabbitChannel) error
	DeclareExchange(ctx context.Context, exchange Exchange) error
}

func NewProducer(parameters ConnectionParameters) (Producer, error) {
	client := newClient(parameters)

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return &producer{
		client,
	}, nil
}

func (p *producer) CreateConfirmedChannel(ctx context.Context) (RabbitChannel, error) {
	theChannel, err := p.client.CreateChannel(ctx)
	if err != nil {
		return nil, err
	}

	if err := theChannel.Confirmation(false); err != nil {
		return nil, err
	}

	return theChannel, nil
}

func (p *producer) ProduceWithChannel(ctx context.Context, routingKey, exchangeName string, payload interface{}, priority uint8, channel RabbitChannel) error {
	data, err := prepareData(payload, amqp.Table{}, priority)
	if err != nil {
		return err
	}

	if channel == nil {
		return fmt.Errorf("channel should not be nil")
	}

	if err := channel.PublishMessage(ctx, exchangeName, routingKey, data); err != nil {
		return err
	}

	return nil
}

func (p *producer) DeclareExchange(ctx context.Context, exchange Exchange) error {
	theChannel, err := p.client.CreateChannel(ctx)
	defer func() { _ = theChannel.Close(ctx) }()
	if err != nil {
		return err
	}
	return theChannel.DeclareExchange(exchange)
}

func prepareData(payload interface{}, headers amqp.Table, priority uint8) (amqp.Publishing, error) {
	var (
		data []byte
		err  error
	)

	switch v := payload.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	case nil:
		data = nil
	default:
		data, err = json.Marshal(payload)

		if err != nil {
			return amqp.Publishing{}, fmt.Errorf("failed to convert data, err:%v", err)
		}
	}

	return amqp.Publishing{
		Headers:         headers,
		ContentType:     "application/json",
		ContentEncoding: "UTF-8",
		Body:            data,
		DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
		Priority:        priority,       // 0-9
	}, nil

}
