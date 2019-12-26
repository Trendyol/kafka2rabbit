package kafka_test

import (
	"github.com/kafka2rabbit/internal/integrationcontainers/kafka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestProducer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kafka Suite")
}

var (
	KafkaContainer *kafka.Container
)

var _ = BeforeSuite(func() {
	KafkaContainer = kafka.NewContainer("johnnypark/kafka-zookeeper")
	KafkaContainer.Run()
})

var _ = AfterSuite(func() {
	KafkaContainer.ForceRemoveAndPrune()
})
