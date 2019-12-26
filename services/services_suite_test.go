package services_test

import (
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/internal/integrationcontainers/kafka"
	"github.com/kafka2rabbit/services/event_executor"
	"oms-gitlab.trendyol.com/oms/event-generator/appconfig"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var (
	KafkaContainer *kafka.Container
)

var _ = BeforeSuite(func() {
	KafkaContainer = kafka.NewContainer(appconfig.KAFKA_IMAGE)
	KafkaContainer.Run()

})

var _ = AfterSuite(func() {
	KafkaContainer.ForceRemoveAndPrune()
})

type mockEventExecutor struct{}

func NewMockEventExecutor() event_executor.Executor {
	return &mockEventExecutor{}
}

func (m *mockEventExecutor) Execute(message *sarama.ConsumerMessage) error {
	return nil
}
