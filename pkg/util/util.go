package util

import (
	"github.com/Shopify/sarama"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/responses"
	"github.com/kafka2rabbit/pkg/kafka"
	. "github.com/onsi/ginkgo"
	"net/http"
	"sync"
	"time"
)

func WaitUntil(predicate func() bool, timeout time.Duration, interval time.Duration) bool {
	var totalWait time.Duration

	for predicate() == false && totalWait < timeout {
		time.Sleep(interval)
		totalWait = totalWait + interval
	}

	return totalWait < timeout
}

var BeforeAll = func(beforeAllFunc func()) {
	var once sync.Once

	BeforeEach(func() {
		once.Do(func() {
			beforeAllFunc()
		})
	})
}

var AfterAll = func(afterAllFunc func()) {
	var once sync.Once

	AfterEach(func() {
		once.Do(func() {
			afterAllFunc()
		})
	})
}

func Consume(kafkaConfig kafka.ConnectionParameters) (string, string, []*sarama.RecordHeader) {
	v, _ := sarama.ParseKafkaVersion(kafkaConfig.Version)

	config := sarama.NewConfig()
	config.Version = v
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer([]string{kafkaConfig.Brokers}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	partitions, _ := master.Partitions(kafkaConfig.Topic)
	partitionConsumer, err := master.ConsumePartition(kafkaConfig.Topic, partitions[0], sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	msg := <-partitionConsumer.Messages()
	return string(msg.Key), string(msg.Value), msg.Headers
}

func NotFoundStub() gobank.StubElement {
	return gobank.Stub().
		Responses(responses.
			Is().
			StatusCode(http.StatusNotFound).
			Header("Content-Type", "application/json").
			Header("X-Mountebank-Catch-All", "Matched").
			Build()).
		Build()
}