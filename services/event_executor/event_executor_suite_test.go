package event_executor_test

import (
	"github.com/durmaze/gobank"
	"github.com/kafka2rabbit/internal/integrationcontainers/kafka"
	"github.com/kafka2rabbit/internal/integrationcontainers/mountebank"
	"github.com/kafka2rabbit/internal/integrationcontainers/rabbit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEventExecutor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EventExecutor Suite")
}

var (
	kafkaContainer      *kafka.Container
	rabbitContainer     *rabbit.Container
	mountebankContainer *mountebank.Container
	Mountebank          *gobank.Client
)

var _ = BeforeSuite(func() {
	kafkaContainer = kafka.NewContainer("johnnypark/kafka-zookeeper")
	kafkaContainer.Run()
	rabbitContainer = rabbit.NewContainer("rabbitmq:3.7.17-management")
	rabbitContainer.Run()
	mountebankContainer = mountebank.NewContainer("andyrbell/mountebank:2.1.0")
	mountebankContainer.Run()

	Mountebank = gobank.NewClient(mountebankContainer.AdminUri())
})

var _ = AfterSuite(func() {
	kafkaContainer.ForceRemoveAndPrune()
})
