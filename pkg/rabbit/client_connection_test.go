package rabbit_test

import (
	"github.com/kafka2rabbit/pkg/rabbit"
	. "github.com/kafka2rabbit/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When RabbitMQ client connects", func() {
	Context("with valid parameters", func() {
		var (
			connectionParameters rabbit.ConnectionParameters
			connectionErr        error
		)

		BeforeAll(func() {
			connectionParameters = rabbit.ConnectionParameters{
				ConnectionString: rabbitMQContainer.Address(),
				PrefetchCount:    2,
				RetryCount:       5,
				RetryInterval:    10,
			}

			connectionErr = rabbit.NewClient(connectionParameters).Connect()
		})

		It("should not return error", func() {
			Expect(connectionErr).NotTo(HaveOccurred())
		})
	})

	Context("with invalid parameters", func() {
		var (
			connectionParameters rabbit.ConnectionParameters
			connectionErr        error
		)

		BeforeAll(func() {
			connectionParameters = rabbit.ConnectionParameters{
				ConnectionString: "amqp://invalid_connection_String",
				PrefetchCount:    2,
				RetryCount:       5,
				RetryInterval:    10,
			}

			var client = rabbit.NewClient(connectionParameters)

			connectionErr = client.Connect()
		})

		It("should return error", func() {
			Expect(connectionErr).Should(Not(BeNil()))
		})
	})
})
