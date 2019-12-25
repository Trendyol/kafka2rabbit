package rabbit

import (
	"context"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

const CorrelationIdHeaderKey = "Correlation-Id"

type ConnectionParameters struct {
	ConnectionString string
	PrefetchCount    int
	RetryCount       int
	RetryInterval    time.Duration
}

type client struct {
	connection           *amqp.Connection
	connectionString     string
	prefetchCount        int
	retryCount           int
	retryInterval        time.Duration
	connectionNotifierCh chan bool
}

type RabbitClient interface {
	Close() error
	Connect() error
	NotifyChannel() chan bool
	CreateChannel(ctx context.Context) (RabbitChannel, error)
}

func newClient(params ConnectionParameters) RabbitClient {
	var rabbitClient = client{
		connectionString:     params.ConnectionString,
		prefetchCount:        params.PrefetchCount,
		retryCount:           params.RetryCount,
		retryInterval:        params.RetryInterval,
		connectionNotifierCh: make(chan bool),
	}

	return &rabbitClient
}

func (c *client) CreateChannel(ctx context.Context) (RabbitChannel, error) {
	return newChannel(ctx, c.connection, c.prefetchCount)
}

func (c *client) Close() error {
	var err error
	if err = c.connection.Close(); err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (c *client) Connect() error {
	var err error

	retryCount := 0
	for retryCount != c.retryCount {
		if retryCount != 0 {
			time.Sleep(c.retryInterval)
		}

		retryCount++
		if c.connection, err = amqp.Dial(c.connectionString); err != nil {
			c.notifyConnectionStatus(false)
			continue
		}

		c.notifyConnectionStatus(true)

		c.onClose()

		break
	}

	if retryCount == c.retryCount {
		return fmt.Errorf("connection error")
	}

	return nil
}

func (c *client) NotifyChannel() chan bool {
	return c.connectionNotifierCh
}

func (c *client) notifyConnectionStatus(status bool) {
	go func() { c.connectionNotifierCh <- status }()
}

func (c *client) onClose() {
	go func() {
		err := <-c.connection.NotifyClose(make(chan *amqp.Error))
		if err == nil {
			return
		}
		c.notifyConnectionStatus(false)
		c.Connect()
	}()
}
