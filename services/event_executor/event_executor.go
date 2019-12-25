package event_executor

import (
	"github.com/Shopify/sarama"
)

type Executor interface {
	Execute(message *sarama.ConsumerMessage) error
}

type EventExecutor struct {
	executor Executor
}

func NewEventExecutor() *EventExecutor {
	return &EventExecutor{}
}

func (c *EventExecutor) SetStrategy(executor Executor) {
	c.executor = executor
}

func (c *EventExecutor) Execute(message *sarama.ConsumerMessage) error {
	return c.executor.Execute(message)
}
