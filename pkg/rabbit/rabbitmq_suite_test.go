package rabbit_test

import (
	"github.com/kafka2rabbit/internal/integrationcontainers/rabbit"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRabbit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RabbitMQ Test Suite")
}

var (
	rabbitMQContainer *rabbit.Container
)

var _ = BeforeSuite(func() {
	rabbitMQContainer = rabbit.NewContainer("rabbitmq:3.7.17-management")
	rabbitMQContainer.Run()

})

var _ = AfterSuite(func() {
	rabbitMQContainer.ForceRemoveAndPrune()
})
