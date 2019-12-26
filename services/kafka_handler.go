package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/pkg/kafka"
	"github.com/kafka2rabbit/services/event_executor"
	"strings"
)

type eventHandler struct {
	retryExecutor  event_executor.Executor
	errorExecutor  event_executor.Executor
	normalExecutor event_executor.Executor
}

func NewEventHandler(
	retryExecutor event_executor.Executor,
	normalExecutor event_executor.Executor,
	errorExecutor event_executor.Executor) kafka.EventHandler {
	return &eventHandler{
		retryExecutor:  retryExecutor,
		normalExecutor: normalExecutor,
		errorExecutor:  errorExecutor,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (e *eventHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("kafka 2 rabbit listener is starting")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (e *eventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (e *eventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	executor := event_executor.NewEventExecutor()
	e.decideBehavioral(executor, claim)
	for message := range claim.Messages() {
		fmt.Printf("Received messages , msg:%v", string(message.Value))
		err := executor.Execute(message)
		if err != nil {
			fmt.Printf("Error executing message: %+v , err: %+v", message.Value, err)
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func (e *eventHandler) decideBehavioral(eventExecutor *event_executor.EventExecutor, claim sarama.ConsumerGroupClaim) {
	if isRetryTopic(claim) {
		eventExecutor.SetStrategy(e.retryExecutor)
	} else if isErrorTopic(claim) {
		eventExecutor.SetStrategy(e.errorExecutor)
	} else {
		eventExecutor.SetStrategy(e.normalExecutor)
	}
}

func isRetryTopic(claim sarama.ConsumerGroupClaim) bool {
	return strings.Contains(claim.Topic(), "RETRY")
}

func isErrorTopic(claim sarama.ConsumerGroupClaim) bool {
	return strings.Contains(claim.Topic(), "ERROR")
}
